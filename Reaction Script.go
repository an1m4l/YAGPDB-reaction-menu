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
