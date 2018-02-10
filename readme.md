# revelforce-admin

## Config
`config.[env-name].env.json`

**env** will come from environment variables

## Flash
**add flash message**
```
err := flash.Add(w, r, "it's me", "success")
if err != nil {
	serverError(w, r, err)
	return
}
```