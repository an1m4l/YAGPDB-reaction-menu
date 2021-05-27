{{/*trigger script (can be used multiple times), use trigger type: Command (mention/cmd prefix)*/}}
{{define "settings"}}
{{/*emoji formats a:name:id | name:id | name (unicode character)*/}}
{{.Set "defemoji" (cslice "1️⃣" "2️⃣" "3️⃣" "4️⃣" "5️⃣" "6️⃣" "7️⃣" "8️⃣" "9️⃣" "🔟")}}{{/*Default emoji when not set in the embeds (you can use the same order you want for all and save space with this)*/}}
{{.Set "closeemoji" "❌"}}{{/*set the emoji for closing the embed (this isn't included in the public one)*/}}
{{.Set "publicarg" (cslice 00000000000000)}}{{/*set the role id's who can make a embed menu public and doesn't get deleted after the duration*/}}
{{.Set "mentionarg" (cslice 000000000000000)}}{{/*set the role id's who can make a embed for a user with a mention*/}}
{{.Set "duration" 300}}{{/*set the duration in seconds (max 24h) when an embed auto gets deleted*/}}
{{/*.Set "channelid" 0}}{{/*this set the default output channel, for trigger channel delete this line or comment out (default)*/}}

{{/*first embed is your start page (menu)*/}}
{{$embeds:= cslice
	(sdict
		"emoji" "a:yagseesall:598726566075301928"
		"title" "Page 1"
		"description" "page 1 description: All embed can be used inside the menu parts"
	)
	(sdict
		"title" "Page 2"
		"description" "page 2 description"
	)
	(sdict
		"title" "Page 3"
		"description" "page 3 description"
	)
	(sdict
		"title" "Page 4"
		"description" "page 4 description"
	)
}}
{{.Set "embeds" $embeds}}
{{end}}

{{define "checkrole"}}
    {{range .mr}}{{if not $.ret}}{{if in $.roles .}}{{$.Set "ret" 1}}{{end}}{{end}}{{end}}
{{end}}

{{template "settings" ($inst:= sdict)}}
{{$inst.Set "users" (cslice .User.ID)}}
{{if .CmdArgs}}
    {{$perms:= 1}}
    {{if eq (lower (index .CmdArgs 0)) "public"}}
        {{if $inst.publicarg}}{{template "checkrole" ($x:= sdict "mr" .Member.Roles "roles" $inst.publicarg)}}{{if not $x.ret}}{{$perms = 0}}{{end}}{{end}}
        {{if $perms}}{{with $inst}}{{.Set "duration" 0}}{{.Del "channelid"}}{{.Del "closeemoji"}}{{.Del "users"}}{{end}}{{end}}
    {{else if ($mention:= userArg (index .CmdArgs 0))}}
        {{if $inst.mentionarg}}{{template "checkrole" ($x:= sdict "mr" .Member.Roles "roles" $inst.mentionarg)}}{{if not $x.ret}}{{$perms = 0}}{{end}}{{end}}
        {{if $perms}}{{if not (in $inst.users $mention.ID)}}{{$inst.Set "users" (($inst.Get "users").Append $mention.ID)}}{{end}}{{end}}
    {{end}}
{{end}}

{{$msgid:= sendMessageRetID nil (cembed (index $inst.embeds 0))}}
{{range $i, $e:= $inst.embeds}}{{if not .emoji}}{{.Set "emoji" (index $inst.defemoji $i)}}{{end}}{{if not (eq $i 0)}}{{addMessageReactions $inst.channelid $msgid .emoji}}{{end}}{{end}}
{{if $inst.closeemoji}}{{addMessageReactions $inst.channelid $msgid $inst.closeemoji}}{{end}}
{{if $inst.duration}}
    {{deleteMessage $inst.channelid $msgid $inst.duration}}
    {{dbSetExpire 0 (print "EM" .Channel.ID $msgid) $inst $inst.duration}}
{{else}}
    {{dbSet 0 (print "EM" .Channel.ID $msgid) $inst}}
{{end}}
{{deleteTrigger}}