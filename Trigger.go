{{/*trigger script (can be used multiple times), use trigger type: Command (mention/cmd prefix)*/}}
{{define "settings"}}
{{/*emoji formats a:name:id | name:id | name (unicode character)*/}}
{{.Set "defemoji" (cslice "1️⃣" "2️⃣" "3️⃣" "4️⃣" "5️⃣" "6️⃣" "7️⃣" "8️⃣" "9️⃣" "🔟")}}{{/*Default emoji when not set in the embeds (you can use the same order you want for all and save space with this)*/}}
{{.Set "closeemoji" "❌"}}{{/*set the emoji for closing the embed (this isn't included in the public one)*/}}
{{.Set "publicarg" (cslice 00000000000000)}}{{/*set the rolesid who can make a embed menu public*/}}
{{.Set "mentionarg" (cslice 000000000000000)}}{{/*set the rolesid who can make a embed for a user with a mention*/}}
{{.Set "duration" 300}}{{/*set the duration in seconds when an embed auto gets deleted*/}}
{{/*.Set "channelid" 0}}{{/*this set the output channel, for trigger channel delete this line*/}}

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
{{deleteTrigger}}{{/*reaction script, use trigger type: Reaction/Added Reactions Only*/}}

{{/*Converter for interfase types from DB*/}}
{{- define "standardize"}}{{/*Thanks to Satty*/}}
{{- $val:= (.Get "val")}}{{$rDict := sdict}}{{$rVal := ""}}
    {{- if (eq (printf "%T" $val) `map[string]interface {}`)}}{{$rVal = sdict $val}}
    {{- else if (eq (printf "%T" $val) `templates.SDict`)}}{{$rVal = $val}}
    {{- else if (eq (printf "%T" $val) `[]interface {}`)}}{{$rVal = cslice.AppendSlice $val}}
    {{- else if (eq (printf "%T" $val) `templates.Slice`)}}{{$rVal = $val}}{{end}}
    {{- if (print $rVal)}}
    {{- range $k,$v := $rVal}}
    {{- if in (cslice `map[string]interface {}` `[]interface {}` `templates.SDict` `templates.Slice`) (printf "%T" $val)}}
    {{- $rDict.Set "val" $v}}{{template "standardize" $rDict}}{{$rVal.Set $k ($rDict.Get "ret")}}{{end}}{{end}}
    {{- else}}{{$rVal = $val}}{{end}}
{{- (.Set "ret" $rVal)}}
{{- end}}

{{define "getEmojiName"}}
	{{with .emoji}}{{$.Set "name" (print (or (and .Animated "a:") "") .Name (or (and .ID (print ":" .ID)) ""))}}{{end}}
{{end}}

{{if .ReactionAdded}}
	{{template "getEmojiName" ($emoji := sdict "emoji" .Reaction.Emoji)}}
	{{if ($db:= (dbGet 0 (print "EM" .Reaction.ChannelID .Reaction.MessageID)).Value)}}
		{{$delReaction:= 1}}
		{{template "standardize" ($x:= sdict "val" $db)}}{{$inst:= $x.ret}}
		{{if or (not $inst.users) (in $inst.users .Reaction.UserID)}}
			{{range $i, $e:= $inst.embeds}}{{if .emoji}}{{if eq .emoji $emoji.name}}{{$emoji.Set "pos" $i}}{{$emoji.Set "match" 1}}{{end}}{{end}}{{end}}
			{{if $emoji.match}}
				{{$delReaction:= 0}}
				{{editMessage .Reaction.ChannelID .Reaction.MessageID (cembed (index $inst.embeds $emoji.pos))}}
				{{deleteAllMessageReactions .Channel.ID .Message.ID}}
				{{range $i, $e:= $inst.embeds}}{{if .emoji}}{{if not (eq $i $emoji.pos)}}{{addMessageReactions $.Reaction.ChannelID $.Reaction.MessageID $e.emoji}}{{end}}{{end}}{{end}}
				{{if $inst.closeemoji}}{{addMessageReactions .Reaction.ChannelID .Reaction.MessageID $inst.closeemoji}}{{end}}
			{{end}}
			{{if $inst.closeemoji}}
				{{if eq $inst.closeemoji $emoji.name}}
					{{$delReaction:= 0}}
					{{dbDel 0 (print "EM" .Reaction.ChannelID .Reaction.MessageID)}}
					{{deleteMessage .Reaction.ChannelID .Reaction.MessageID 1}}
				{{end}}
			{{end}}
		{{end}}
		{{if $delReaction}}
			{{deleteMessageReaction .Reaction.ChannelID .Reaction.MessageID .Reaction.UserID $emoji.name}}
		{{end}}
	{{end}}
{{end}}
