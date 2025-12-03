---
title: Vars
var11: Value from page
---

Var1: `{{ .var1 }}`
Var2: `{{ .var2 }}`
Var3: `{{ .var3 }}`
Var4: `{{ .var4 }}`
Var5: `{{ .var5 }}`
Var6: `{{ .var6 }}`
Var7: `{{ .var7 }}`
Var8: `{{ .var8 }}`
Var9: `{{ .var9 }}`
Var10: `{{ .var10 }}`
Var11: `{{ .var11 }}`

---

System.BuildTime: `{{ .System.BuildTime }}` 
System.Now: `{{ .System.Now }}` 

---

Site.Name: `{{ .Site.Name }}`
Site.BaseUrl: `{{ .Site.BaseUrl }}`
Site.Author.Name: `{{ .Site.Author.Name }}`
Site.Author.Email: `{{ .Site.Author.Email }}`
Site.Author.Website: `{{ .Site.Author.Website }}`
Site.Copyright.Year: `{{ .Site.Copyright.Year }}`
Site.Copyright.Title: `{{ .Site.Copyright.Title }}`
Site.Copyright.Rights: `{{ .Site.Copyright.Rights }}`
Site.Logo.Url: `{{ .Site.Logo.Url }}`

---

Page.Id: `{{ .Page.Id }}`
Page.Name: `{{ .Page.Name }}`
Page.Uri: `{{ .Page.Uri }}`
Page.Url: `{{ .Page.Url }}`
Page.Title: `{{ .Page.Title }}`
Page.IsHidden: `{{ .Page.IsHidden }}`
Page.IsDraft: `{{ .Page.IsDraft }}`
Page.IsSystem: `{{ .Page.IsSystem }}`
