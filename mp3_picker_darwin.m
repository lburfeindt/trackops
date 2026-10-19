#import <Cocoa/Cocoa.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>
#import <dispatch/dispatch.h>
#import <stdlib.h>
#import <string.h>

// Wails v2 has separate file and folder dialogs; this panel accepts both.
char *chooseMP3Inputs(void) {
    __block char *selection = NULL;
    dispatch_semaphore_t completed = dispatch_semaphore_create(0);

    dispatch_async(dispatch_get_main_queue(), ^{
        NSOpenPanel *panel = [NSOpenPanel openPanel];
        panel.title = @"Choose MP3 files or folders to sanitize";
        panel.canChooseFiles = YES;
        panel.canChooseDirectories = YES;
        panel.allowsMultipleSelection = YES;
        panel.resolvesAliases = NO;
        UTType *mp3Type = [UTType typeWithFilenameExtension:@"mp3"];
        if (mp3Type != nil) {
            panel.allowedContentTypes = @[mp3Type];
        }

        NSWindow *window = [NSApp mainWindow];
        if (window == nil) {
            selection = NULL;
            dispatch_semaphore_signal(completed);
            return;
        }

        [panel beginSheetModalForWindow:window completionHandler:^(NSModalResponse response) {
            NSMutableArray *paths = [NSMutableArray array];
            if (response == NSModalResponseOK) {
                for (NSURL *url in panel.URLs) {
                    [paths addObject:url.path];
                }
            }
            NSData *data = [NSJSONSerialization dataWithJSONObject:paths options:0 error:NULL];
            NSString *json = [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease];
            selection = strdup(json.UTF8String);
            dispatch_semaphore_signal(completed);
        }];
    });

    dispatch_semaphore_wait(completed, DISPATCH_TIME_FOREVER);
    dispatch_release(completed);
    return selection;
}
