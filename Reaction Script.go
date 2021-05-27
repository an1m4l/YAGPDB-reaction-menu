{{/*reaction script, use trigger type: Reaction/Added Reactions Only*/}}

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
	{{$delReaction:= 1}}
	{{template "getEmojiName" ($emoji := sdict "emoji" .Reaction.Emoji)}}
	{{$reaction:=(print (or (and ($x:=($z:=.Reaction.Emoji).Animated) "a:") "") $z.Name (or (and $z.ID (print ":" $z.ID)) ""))}}
	{{if ($db:= (dbGet 0 (print "EM" .Reaction.ChannelID .Reaction.MessageID)).Value)}}
		{{template "standardize" ($x:= sdict "val" $db)}}{{$inst:= $x.ret}}
		{{if or (not $inst.users) (in $inst.users .Reaction.UserID)}}
			{{$match:= 0}}
			{{$pos:= 0}}
			{{range $i, $e:= $inst.embeds}}{{if .emoji}}{{if eq .emoji $reaction}}{{$pos = $i}}{{$match = 1}}{{end}}{{end}}{{end}}
			{{if $match}}
				{{$delReaction:= 0}}
				{{editMessage .Reaction.ChannelID .Reaction.MessageID (cembed (index $inst.embeds $pos))}}
				{{deleteAllMessageReactions .Channel.ID .Message.ID}}
				{{range $i, $e:= $inst.embeds}}{{if .emoji}}{{if not (eq $i $pos)}}{{addMessageReactions $.Reaction.ChannelID $.Reaction.MessageID $e.emoji}}{{end}}{{end}}{{end}}
				{{if $inst.closeemoji}}{{addMessageReactions .Reaction.ChannelID .Reaction.MessageID $inst.closeemoji}}{{end}}
			{{end}}
			{{if $inst.closeemoji}}
				{{if eq $inst.closeemoji $reaction}}
					{{$delReaction:= 0}}
					{{dbDel 0 (print "EM" .Reaction.ChannelID .Reaction.MessageID)}}
					{{deleteMessage .Reaction.ChannelID .Reaction.MessageID 1}}
				{{end}}
			{{end}}
		{{end}}
	{{end}}
	{{if $delReaction}}
		{{deleteMessageReaction .Reaction.ChannelID .Reaction.MessageID .Reaction.UserID $reaction}}
	{{end}}
{{end}}
