package utils

import (
	"net/http"
	apiModel "openstreetmap-go/src/api/model"
	"openstreetmap-go/src/utils"
	"strconv"
	"strings"
	"time"
)

func LogAPIAccess(
	w *apiModel.WrappedWriter,
	r *http.Request,
	startTime time.Time,
	tag string,
	messages []string,
) {
	//var ctx, _ = apiUtils.GetWrappedContext(r)

	//var jwtPayload = ctx.TagAlong.Data

	var elapsedTime = time.Since(startTime).Milliseconds()

	var latencyIsBad = elapsedTime > 2000
	var latencyIsGood = elapsedTime <= 900

	var elapsedTimeColor = utils.Yellow

	if latencyIsBad {
		elapsedTimeColor = utils.Red
	} else if latencyIsGood {
		elapsedTimeColor = utils.Cyan
	}

	//var userType, hasUserType = jwtPayload.UserType, jwtPayload.UserType != ""
	//var userId, hasUserId = jwtPayload.UserId, jwtPayload.UserId != ""
	//var sessionId, hasSessionId = jwtPayload.SessionId, jwtPayload.SessionId != ""

	//var sessionIdColor = utils.Cyan
	//var userIdColor = utils.Yellow
	//var userTypeColor = utils.LightBlue

	//if !hasUserType || userType == "" {
	//	userType = "unknown"
	//	userTypeColor = utils.BrightRed
	//}
	//
	//if !hasUserId || userId == "" {
	//	userId = "unrecognized"
	//	userIdColor = utils.BrightRed
	//}
	//
	//if !hasSessionId || sessionId == "" {
	//	sessionId = "unauthenticated"
	//	sessionIdColor = utils.BrightRed
	//}

	var responseIsServerError = w.StatusCode >= 500
	var responseIsUserError = w.StatusCode >= 400
	var responseIsRedirect = w.StatusCode >= 300
	var responseIsSuccess = w.StatusCode >= 200

	var statusColor = utils.Cyan

	if responseIsServerError {
		statusColor = utils.BrightRed
	} else if responseIsUserError {
		statusColor = utils.Red
	} else if responseIsRedirect {
		statusColor = utils.Yellow
	} else if responseIsSuccess {
		statusColor = utils.Green
	}

	var logTimeSize = "                    "

	var message = ""

	var indentation = logTimeSize + utils.GetEquivalentWhiteSpace(tag) + " "

	for index, part := range messages {
		if index == 0 {
			message = message + part
		} else {
			message = message + " " + indentation
		}

		message = message + "\n"
	}

	if len(messages) > 0 {
		message = message + "\n" + indentation + "|" + "\n"
		message = message + indentation
	}

	message = message + utils.MagentaString(strings.ToUpper(r.Method)) + " "
	message = message + r.URL.Path + " "
	message = message + statusColor + utils.ColorString(strconv.FormatInt(int64(w.StatusCode), 10), statusColor) + "\n"

	//message = message + indentation + utils.ColorString(userType, userTypeColor) + " "
	//message = message + utils.ColorString(userId, userIdColor) + " "
	//message = message + utils.ColorString(sessionId, sessionIdColor) + "\n"

	message = message + indentation + "from " + utils.MagentaString(GetRequestIP(r)) + "\n"
	message = message + indentation + "via " + utils.GreyString(r.Header.Get("user-agent")) + "\n"
	message = message + indentation + "on " + utils.GreyString(time.Now().UTC().String()) + "\n"
	message = message + indentation + "took " + utils.ColorString(strconv.FormatInt(elapsedTime, 10), elapsedTimeColor) + " ms"

	if responseIsUserError || responseIsServerError {
		utils.LogError(tag, message)
	} else if responseIsRedirect {
		utils.LogWarn(tag, message)
	} else if responseIsSuccess {
		utils.LogSuccess(tag, message)
	}
}
