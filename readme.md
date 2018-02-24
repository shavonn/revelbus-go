# RevelForce

Revel Bus backend, tour management.

## Installation

###  Config
`config/config.[env-name].env.json`
**env** environment variable

## Usage

### Flash
Add flash message
```
err := flash.Add(w, r, "it's me", "success")
if err != nil {
	view.ServerError(w, r, err)
	return
}
```