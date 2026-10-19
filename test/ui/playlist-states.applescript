on run arguments
	-- Given
	set firstPlaylist to item 1 of arguments
	set secondPlaylist to item 2 of arguments
	set emptyPlaylist to item 3 of arguments
	set invalidPlaylist to item 4 of arguments
	tell application "TrackOps" to activate

	-- When
	my choosePlaylist(firstPlaylist)

	-- Then
	my waitForText("First Artist - Café")
	my waitForText("Second Artist, Club Mix")
	tell application "System Events" to tell process "TrackOps"
		set copyButton to my findButton("Copy titles")
		if not enabled of copyButton then error "Copy titles was disabled for a populated playlist"
		click copyButton
	end tell

	-- When
	my openPicker()
	tell application "System Events" to tell process "TrackOps"
		key code 53
	end tell

	-- Then
	my waitForText("First Artist - Café")

	-- When
	my choosePlaylist(secondPlaylist)

	-- Then
	my waitForText("Replacement Title")
	my assertTextAbsent("First Artist - Café")

	-- When
	my choosePlaylist(emptyPlaylist)

	-- Then
	my waitForText("No titles found")
	my assertTextAbsent("Replacement Title")
	tell application "System Events" to tell process "TrackOps"
		if enabled of (my findButton("Copy titles")) then error "Copy titles was enabled for an empty playlist"
	end tell

	-- When
	my choosePlaylist(invalidPlaylist)

	-- Then
	my waitForTextContaining("token too long")
	my assertTextAbsent("No titles found")
	tell application "System Events" to tell process "TrackOps"
		if enabled of (my findButton("Copy titles")) then error "Copy titles was enabled after a read error"
		set windowElements to get entire contents of window "TrackOps"
		repeat with element in windowElements
			if role of element is "AXButton" and name of element is "Save titles" then error "Save titles was offered after a read error"
		end repeat
	end tell

	-- When
	my choosePlaylist(firstPlaylist)

	-- Then
	my waitForText("First Artist - Café")
end run

on findButton(buttonName)
	tell application "System Events" to tell process "TrackOps"
		repeat 30 times
			set windowElements to get entire contents of window "TrackOps"
			repeat with element in windowElements
				if role of element is "AXButton" and name of element is buttonName then return element
			end repeat
			delay 0.2
		end repeat
	end tell
	error buttonName & " button is missing"
end findButton

on openPicker()
	tell application "System Events" to tell process "TrackOps"
		click my findButton("Format Playlist")
		repeat 30 times
			if exists sheet 1 of window "TrackOps" then return
			delay 0.2
		end repeat
	end tell
	error "Native playlist picker did not open"
end openPicker

on choosePlaylist(playlistPath)
	my openPicker()
	tell application "System Events" to tell process "TrackOps"
		set frontmost to true
		keystroke "g" using {command down, shift down}
		repeat 30 times
			if exists sheet 1 of sheet 1 of window "TrackOps" then exit repeat
			delay 0.2
		end repeat
		set value of text field 1 of sheet 1 of sheet 1 of window "TrackOps" to playlistPath
		key code 36
		repeat 30 times
			try
				if enabled of button "Open" of splitter group 1 of sheet 1 of window "TrackOps" then exit repeat
			end try
			delay 0.2
		end repeat
		click button "Open" of splitter group 1 of sheet 1 of window "TrackOps"
		repeat 30 times
			if not exists sheet 1 of window "TrackOps" then return
			delay 0.2
		end repeat
	end tell
	error "Native playlist picker did not close"
end choosePlaylist

on waitForText(expectedText)
	repeat 30 times
		if my hasText(expectedText) then return
		delay 0.2
	end repeat
	error expectedText & " was not visible"
end waitForText

on waitForTextContaining(expectedText)
	repeat 30 times
		tell application "System Events" to tell process "TrackOps"
			set windowElements to get entire contents of window "TrackOps"
			repeat with element in windowElements
				try
					if (value of element as text) contains expectedText then return
				end try
			end repeat
		end tell
		delay 0.2
	end repeat
	error expectedText & " was not visible"
end waitForTextContaining

on assertTextAbsent(unexpectedText)
	if my hasText(unexpectedText) then error unexpectedText & " remained visible"
end assertTextAbsent

on hasText(expectedText)
	tell application "System Events" to tell process "TrackOps"
		set windowElements to get entire contents of window "TrackOps"
		repeat with element in windowElements
			try
				if value of element as text is expectedText then return true
			end try
		end repeat
	end tell
	return false
end hasText
