{{/*reaction script, use trigger type: Reaction/Added Reactions Only*/}}

{{/*Converter for interfase types from DB*/}}
{{- define "standardize"}}
    {{- $types:=sdict `templates.SDict` 1 `templates.Slice` 1 `map[string]interface {}` 2 `[]interface {}` 3}}
    {{- $rVal:= ""}}
    {{- if $t:= $types.Get (printf "%T" .val)}}
        {{- if eq $t 1}}{{$rVal = .val}}{{else if eq $t 2}}{{$rVal = sdict .val}}{{else if eq $t 3}}{{$rVal = cslice.AppendSlice .val}}{{end}}
        {{- range $k,$v := $rVal}}{{if $types.Get (printf "%T" $v)}}{{template "standardize" $x:= sdict "val" $v}}{{$rVal.Set $k $x.ret}}{{end}}{{end}}
    {{- end}}
    {{- .Set "ret" (or $rVal .val)}}
{{- end}}

{{if .ReactionAdded}}
	{{$delReaction:= 1}}
	{{$reaction:=(print (or (and ($x:=($z:=.Reaction.Emoji).Animated) "a:") "") $z.Name (or (and $z.ID (print ":" $z.ID)) ""))}}
	{{if ($inst:= (dbGet 0 (print "EM" ($cID:= .Reaction.ChannelID) ($mID:= .Reaction.MessageID))).Value)}}
		{{/*template "standardize" ($x:= sdict "val" $db)}}{{$inst:= $x.ret}}
		{{$db}}
		{{(printf "%T" $db)}}
		{{(printf "%T" (index $db.defemoji 0))}}
		{{(printf "%T" $db.duration)}}
		{{(printf "%T" (index $db.embeds 0))}}
		{{(printf "%T" $db.mentionarg)}}
		{{(printf "%T" $db.publicargarg)}}
		{{(printf "%T" (index $db.users 0))*/}}
		```{{(printf "%T" .Reaction.UserID)}}```
		```{{(printf "%T" (index $inst.users 0))}}```
		{{if or (not $inst.users) (in $inst.users .Reaction.UserID)}}in users
			{{$match:= 0}}
			{{$pos:= 0}}
			{{range $i, $e:= $inst.embeds}}{{if .emoji}}{{if eq .emoji $reaction}}{{$pos = $i}}{{$match = 1}}{{end}}{{end}}{{end}}
			{{if $match}}
				{{$delReaction:= 0}}
				{{editMessage $cID $mID (cembed (index $inst.embeds $pos))}}
				{{deleteAllMessageReactions .Channel.ID .Message.ID}}
				{{range $i, $e:= $inst.embeds}}{{if .emoji}}{{if not (eq $i $pos)}}{{addMessageReactions $cID $mID $e.emoji}}{{end}}{{end}}{{end}}
				{{if $inst.closeemoji}}{{addMessageReactions $cID $mID $inst.closeemoji}}{{end}}
			{{end}}
			{{if $inst.closeemoji}}
				{{if eq $inst.closeemoji $reaction}}
					{{$delReaction:= 0}}
					{{dbDel 0 (print "EM" $cID $mID)}}
					{{deleteMessage $cID $mID 1}}
				{{end}}
			{{end}}
		{{end}}
	{{end}}
	{{if $delReaction}}
		{{deleteMessageReaction .Reaction.ChannelID .Reaction.MessageID .Reaction.UserID $reaction}}
	{{end}}
{{end}}