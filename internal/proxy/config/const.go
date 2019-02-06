package config

const configVersionVar = "$(version)"
const configMethodVar = "$(method)"
const configHostVar = "$(host)"
const configPathVar = "$(path)"
const configContentTypeVar = "$(contentType)"
const configStatusCodeVar = "$(statusCode)"
const configReqHeaderVar = "$(req.header::<NAME>)"
const configRespHeaderVar = "$(resp.header::<NAME>)"
const configCookieVar = "$(cookie::<NAME>)"

// EvalVarPrefix ...
const EvalVarPrefix = "Krat"

// EvalVersionVar ...
const EvalVersionVar = EvalVarPrefix + "VERSION"

// EvalMethodVar ...
const EvalMethodVar = EvalVarPrefix + "METHOD"

// EvalHostVar ...
const EvalHostVar = EvalVarPrefix + "HOST"

// EvalPathVar ...
const EvalPathVar = EvalVarPrefix + "PATH"

// EvalContentTypeVar ...
const EvalContentTypeVar = EvalVarPrefix + "CONTENTTYPE"

// EvalStatusCodeVar ...
const EvalStatusCodeVar = EvalVarPrefix + "STATUSCODE"

// EvalReqHeaderVar ...
const EvalReqHeaderVar = EvalVarPrefix + "REQHEADER"

// EvalRespHeaderVar ...
const EvalRespHeaderVar = EvalVarPrefix + "RESPHEADER"

// EvalCookieVar ...
const EvalCookieVar = EvalVarPrefix + "COOKIE"
