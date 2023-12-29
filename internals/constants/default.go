package constants

import "time"

const LengthOfID = 16

const DefaultPage = 1
const DefaultSize = 10
const MaximumSize = 1e5

const DefaultTimeout = 15 * time.Second

const AuthorizationHeader = "Authorization"
const AuthorizationKey = "Bearer"
