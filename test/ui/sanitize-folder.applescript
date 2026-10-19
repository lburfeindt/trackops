on run arguments
	-- Given
	set folderPath to item 1 of arguments
	set emptyFolderPath to item 2 of arguments
	set stopFolderPath to item 3 of arguments
	set stopStatusPath to item 4 of arguments
	tell application "TrackOps" to activate

	-- When
	my chooseFolder(folderPath)

	-- Then
	my waitForText("Updated: 2 · Unchanged: 0 · Failed: 1")
	my waitForTextContaining("broken.mp3")

	-- When
	my chooseFolder(emptyFolderPath)

	-- Then
	my waitForText("No MP3 files found.")

	-- When
	my chooseFolder(stopFolderPath)

	-- Then
	my waitForTextContaining("Sanitizing: 0002.mp3")
	my clickButton("Stop")
	set resultText to my waitForTextContaining("Stopped · Updated:")
	do shell script "printf %s " & quoted form of resultText & " > " & quoted form of stopStatusPath
end run

on clickButton(buttonName)
	tell application "System Events" to tell process "TrackOps"
		click my findButton(buttonName)
	end tell
end clickButton

on findButton(buttonName)
	tell application "System Events" to tell process "TrackOps"
		repeat 50 times
			set windowElements to get entire contents of window "TrackOps"
			repeat with element in windowElements
				if role of element is "AXButton" and name of element is buttonName then return element
			end repeat
			delay 0.2
		end repeat
	end tell
	error buttonName & " button is missing"
end findButton

on chooseFolder(folderPath)
	my clickButton("MP3 folder")
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
		set value of text field 1 of sheet 1 of sheet 1 of window "TrackOps" to folderPath
		key code 36
		repeat 30 times
			if not exists sheet 1 of sheet 1 of window "TrackOps" then exit repeat
			delay 0.2
		end repeat
		click button "Open" of splitter group 1 of sheet 1 of window "TrackOps"
	end tell
end chooseFolder

on waitForText(expectedText)
	repeat 50 times
		if my hasText(expectedText) then return
		delay 0.2
	end repeat
	error expectedText & " was not visible"
end waitForText

on waitForTextContaining(expectedText)
	repeat 50 times
		tell application "System Events" to tell process "TrackOps"
			set windowElements to get entire contents of window "TrackOps"
			repeat with element in windowElements
				try
					set elementText to value of element as text
					if elementText contains expectedText then return elementText
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
