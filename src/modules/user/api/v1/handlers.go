package v1

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"openstreetmap-go/src/api/utils"
	apiValues "openstreetmap-go/src/api/values"
	"openstreetmap-go/src/config"
	gormModel "openstreetmap-go/src/db/generated/gorm"
	db "openstreetmap-go/src/db/generated/prisma"
	"openstreetmap-go/src/modules/acl/aclRespository"
	"openstreetmap-go/src/modules/user/repository"
	"openstreetmap-go/src/modules/userrole/respository"
	utils2 "openstreetmap-go/src/utils"
	"strconv"
	"strings"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request) error {

	body, err := utils.DecodeBody[CreateUserBody](r)
	if err != nil {
		return err
	}
	if body.Password != body.Passwordconfirmation {
		return errors.New("passwords do not match")
	}

	domains := strings.Split(body.Email, "@")[1]
	afterDot := strings.Split(body.Email, ".")[1]

	domainsArray := []string{domains, afterDot}

	acls, err := aclRespository.AclCheck(domainsArray, []string{"allow_account_creation", "no_account_creation"})
	if err != nil {
		return err
	}

	for _, acl := range acls {
		if acl.K == "no_account_creation" {
			return errors.New("no_account_creation")
		}
	}

	sanitizedDisplayName := SanitizeForRegex(body.DisplayName)

	_, err = repository.GetUserByEmailorDisplayName(strings.ToLower(string(body.Email)), sanitizedDisplayName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var now db.DateTime = db.DateTime(time.Now())
	user := gormModel.Users{}

	user.DisplayName = sanitizedDisplayName
	user.PassCrypt = string(hashedPassword)
	user.Email = strings.ToLower(string(body.Email))
	user.CreationTime = now
	user.DataPublic = config.ServerConfigObject.DefaultDataPublic
	user.Description = ""
	user.HomeLat = utils2.PtrOf(config.ServerConfigObject.DefaultHomeLat)
	user.HomeLon = utils2.PtrOf(config.ServerConfigObject.DefaultHomeLon)
	user.HomeZoom = utils2.PtrOf(config.ServerConfigObject.DefaultHomeZoom)
	user.EmailValid = false
	user.Languages = utils2.PtrOf(config.ServerConfigObject.DefaultLanguages)
	user.Status = "pending"
	user.TermsAgreed = utils2.PtrOf(now)
	user.ConsiderPd = config.ServerConfigObject.DefaultConsiderPd
	user.TermsSeen = true
	user.DescriptionFormat = config.ServerConfigObject.DefaultDescriptionFormat
	user.ChangesetsCount = 0
	user.TracesCount = 0
	user.DiaryEntriesCount = 0
	user.ImageUseGravatar = config.ServerConfigObject.DefaultImageUseGravatar
	user.HomeTile = utils2.PtrOf(int64(0))
	user.TouAgreed = utils2.PtrOf(now)
	user.DiaryCommentsCount = utils2.PtrOf(0)
	user.NoteCommentsCount = utils2.PtrOf(0)
	user.CreationAddress = utils2.PtrOf(getClientIP(r))

	_, err = repository.InsertUser(user)
	if err != nil {
		return err
	}

	if err = utils.EncodeBody(w, user, 200); err != nil {
		return apiValues.ApiErrorInternalServerError
	}
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) error {
	body, err := utils.DecodeBody[LoginUserBody](r)
	if err != nil {
		return err
	}

	sanitedDisplayName := SanitizeForRegex(body.Email)
	user, err := repository.GetUserByEmailorDisplayName(strings.ToLower(string(body.Email)), sanitedDisplayName)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassCrypt), []byte(body.Password))
	if err != nil {
		return err
	}
	if err := utils.EncodeBody(w, user, http.StatusOK); err != nil {
		return apiValues.ApiErrorInternalServerError
	}
	return nil
}
func GrantRole(w http.ResponseWriter, r *http.Request) error {
	uid := r.URL.Query().Get("uid")
	role := r.URL.Query().Get("role")

	if uid == "" || role == "" {
		return errors.New("invalid params")
	}

	intId, err := strconv.Atoi(uid)
	if err != nil {
		return err
	}

	var now db.DateTime = db.DateTime(time.Now())

	// Todo use token to get the granter id
	err = respository.InsertUserRole(db.InnerUserRoles{
		UserID:    db.BigInt(intId),
		Role:      db.UserRoleEnum(role),
		CreatedAt: utils2.PtrOf(now),
		UpdatedAt: utils2.PtrOf(now),
		GranterID: db.BigInt(intId),
	})
	if err != nil {
		return err
	}
	w.WriteHeader(201)

	return nil
}

func Confirm(w http.ResponseWriter, r *http.Request) error {
	confirmString := r.URL.Query().Get("confirm_string")
	displayName := r.PathValue("display_name")
	if strings.TrimSpace(displayName) == "" {
		return errors.New("display_name is required")
	}

	tokenSignature := strings.Split(confirmString, "--")
	h := sha1.New()
	h.Write([]byte(tokenSignature[0] + config.ServerConfigObject.SecretKeyBase))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if tokenSignature[0] != expectedSignature {
		return errors.New("invalid token")
	}

	decoded, err := base64.StdEncoding.DecodeString(tokenSignature[1])
	if err != nil {
		return err
	}

	var data []interface{}
	err = json.Unmarshal(decoded, &data)
	if err != nil {
		return err
	}

	userId := data[0].(float64)
	_, err = repository.UpdateUser(
		map[string]interface{}{
			"status":      "active",
			"email_valid": true,
		},
		map[string]interface{}{
			"id": userId,
		},
	)
	if err != nil {
		return err
	}

	w.WriteHeader(201)
	return nil
}
