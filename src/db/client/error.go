package client

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	commonModel "openstreetmap-go/src/common/model"
	commonValues "openstreetmap-go/src/db/values"
)

func IdentifyDatabaseError(err error) (commonModel.ErrorCode, error) {
	if err == nil {
		return commonValues.ErrorCodeDataUnknown, nil
	}

	var _err *pq.Error

	ok := errors.As(err, &_err)

	var tmp = commonValues.DataErrorUnknown

	if !ok {
		if errors.Is(err, sql.ErrNoRows) {
			return commonValues.ErrorCodeDataNotFound, errors.Join(commonValues.DataErrorNotFound, err)
		}

		return commonModel.ErrorCode(tmp.Error()), errors.Join(tmp, err)
	}

	switch _err.Code {
	case "02000":
		tmp = commonValues.DataErrorNotFound
		break
	case "23514", "22001", "22003", "22007", "22012", "22P02", "44000":
		tmp = commonValues.DataErrorInvalid
		break
	case "23505":
		tmp = commonValues.DataErrorDuplicate
		break
	case "23502":
		tmp = commonValues.DataErrorMissingRequiredField
		break
	case "23503":
		tmp = commonValues.DataErrorReferencedRecordNotFound
		break
	}

	return commonModel.ErrorCode(tmp.Error()), errors.Join(tmp, err)
}
