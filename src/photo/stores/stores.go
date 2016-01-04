package stores

const (
	kTableAccount           = "account"
	KIndexAccountByUsername = "username"

	kTableComment                 = "comment"
	KIndexCommentByPhotoId        = "photo_id"
	KIndexCommentByAccountId      = "account_id"
	KIndexCommentByTags           = "tags"
	KIndexCommentByNotificationId = "notification_id"
	KFilterCommentByIsKnown       = "is_known"

	kTablePhoto             = "photo"
	KIndexPhotoById         = "id"
	KIndexPhotoByAccountId  = "account_id"
	KFilterPhotoByIsPrivate = "is_private"
	KFilterPhotoByAccountId = "account_id"

	kTableFollower             = "follower"
	KIndexFollowerByAccountId  = "account_id"
	KIndexFollowerByFollowerId = "follower_id"
)
