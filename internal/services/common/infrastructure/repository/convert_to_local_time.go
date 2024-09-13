package common

import (
	"context"
	"strconv"
	"time"

	settingsDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
)

func (r *CommonRepository) ConvertToLocalTime(ctx context.Context, t time.Time) time.Time {
	userRepo := r.userClient.Client.UserService.Repository
	branch, err := userRepo.GetBranchByID(ctx, userRepo.GetCurrentBranchId(ctx))
	if err != nil {
		return time.Time{}
	}

	var country *settingsDomain.Country
	if err := r.GetSharedDB().First(&country, "country_name = ?", branch.Country).Error; err != nil {
		return time.Time{}
	}

	offsetString := country.UtcOffset[3:]

	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return time.Time{}
	}

	offsetInSeconds := offset * 3600

	location := time.FixedZone("Fixed", offsetInSeconds)

	localTime := t.In(location)
	return time.Date(localTime.Year(), localTime.Month(), localTime.Day(), localTime.Hour(), localTime.Minute(), 0, 0, time.UTC)
}

func (r *CommonRepository) ConvertToLocalTimeByBranch(ctx context.Context, t time.Time, branchId int32) time.Time {
	userRepo := r.userClient.Client.UserService.Repository
	branch, err := userRepo.GetBranchByID(ctx, branchId)
	if err != nil {
		return time.Time{}
	}

	var country *settingsDomain.Country
	if err := r.GetSharedDB().First(&country, "country_name = ?", branch.Country).Error; err != nil {
		return time.Time{}
	}

	offsetString := country.UtcOffset[3:]

	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return time.Time{}
	}

	offsetInSeconds := offset * 3600

	location := time.FixedZone("Fixed", offsetInSeconds)

	localTime := t.In(location)
	return time.Date(localTime.Year(), localTime.Month(), localTime.Day(), localTime.Hour(), localTime.Minute(), 0, 0, time.UTC)
}
