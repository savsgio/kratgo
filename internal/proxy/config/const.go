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

// EvalVersionVar ...
const EvalVersionVar = "KratVERSION"

// EvalMethodVar ...
const EvalMethodVar = "KratMETHOD"

// EvalHostVar ...
const EvalHostVar = "KratHOST"

// EvalPathVar ...
const EvalPathVar = "KratPATH"

// EvalContentTypeVar ...
const EvalContentTypeVar = "KratCONTENTTYPE"

// EvalStatusCodeVar ...
const EvalStatusCodeVar = "KratSTATUSCODE"

// EvalReqHeaderVar ...
const EvalReqHeaderVar = "KratREQHEADER"

// EvalRespHeaderVar ...
const EvalRespHeaderVar = "KratRESPHEADER"

// EvalCookieVar ...
const EvalCookieVar = "KratCOOKIE"
