package user

type Device struct {
	LastAccessTime       int64  `dynamodbav:"last_access_time"`
	AppVersionCode       int32  `dynamodbav:"app_version_code"`
	AppVersionName       string `dynamodbav:"app_version_name"`
	UserAgent            string `dynamodbav:"user_agent"`
	PushToken            string `dynamodbav:"push_token"`
	PushTokenUpdatedTime int64  `dynamodbav:"push_token_updated_time"`
}
