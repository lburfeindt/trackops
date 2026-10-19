on run arguments
	-- Given
	set newPlaylist to item 1 of arguments
	set emptyPlaylist to item 2 of arguments
	set replacePlaylist to item 3 of arguments
	set cancelPlaylist to item 4 of arguments
	set failurePlaylist to item 5 of arguments
	set newOutput to item 6 of arguments
	set emptyOutput to item 7 of arguments
	tell application "TrackOps" to activate

	-- When
	my choosePlaylist(newPlaylist)

	-- Then
	my waitForText("Artist - Café")
	my waitForText("Not saved")
	my assertFileAbsent(newOutput)
	my clickButton("Save titles")
	my waitForText("Saved")

	-- When
	my choosePlaylist(emptyPlaylist)

	-- Then
	my waitForText("No titles found")
	my assertFileAbsent(emptyOutput)
	my clickButton("Save titles")
	my waitForText("Saved")

	-- When
	my choosePlaylist(replacePlaylist)
	my clickButton("Save titles")

	-- Then
	my waitForText("Replace existing output?")
	my clickButton("Replace")
	my waitForText("Saved")

	-- When
	my choosePlaylist(cancelPlaylist)
	my clickButton("Save titles")
	my waitForText("Replace existing output?")
	my clickButton("Cancel")

	-- Then
	my waitForText("Not saved")
	my waitForText("Unsaved Title")

	-- When
	my choosePlaylist(failurePlaylist)
	my clickButton("Save titles")
	my waitForText("Replace existing output?")
	my clickButton("Replace")

	-- Then
	my waitForTextContaining("permission denied")
	my waitForText("Still Visible")
	tell application "System Events" to tell process "TrackOps"
		if not enabled of (my findButton("Copy titles")) then error "Copy titles was disabled after a failed save"
	end tell
end run

on assertFileAbsent(filePath)
	try
		do shell script "test ! -e " & quoted form of filePath
	on error
		error "Opening playlist created output: " & filePath
	end try
end assertFileAbsent

on clickButton(buttonName)
	tell application "System Events" to tell process "TrackOps"
		click my findButton(buttonName)
	end tell
end clickButton

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

on choosePlaylist(playlistPath)
	my clickButton("Format Playlist")
	tell application "System Events" to tell process "TrackOps"
		repeat 30 times
			if exists sheet 1 of window "TrackOps" then exit repeat
			delay 0.2
		end repeat
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
