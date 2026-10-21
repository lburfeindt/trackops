#import <ApplicationServices/ApplicationServices.h>
#import <CoreServices/CoreServices.h>

static CGColorRef color(CGFloat red, CGFloat green, CGFloat blue, CGFloat alpha) {
    CGFloat components[] = {red / 255, green / 255, blue / 255, alpha / 255};
    CGColorSpaceRef space = CGColorSpaceCreateDeviceRGB();
    CGColorRef result = CGColorCreate(space, components);
    CGColorSpaceRelease(space);
    return result;
}

static void ellipse(CGContextRef context, CGFloat x, CGFloat y, CGFloat radius, bool fill) {
    CGRect bounds = CGRectMake(x - radius, y - radius, radius * 2, radius * 2);
    if (fill) CGContextFillEllipseInRect(context, bounds);
    else CGContextStrokeEllipseInRect(context, bounds);
}

int main(int argc, const char *argv[]) {
    if (argc != 2) return 1;
    const size_t size = 1024;
    CGColorSpaceRef space = CGColorSpaceCreateDeviceRGB();
    CGContextRef context = CGBitmapContextCreate(NULL, size, size, 8, 0, space, (CGBitmapInfo)kCGImageAlphaPremultipliedLast);
    if (!context) return 1;
    CGContextTranslateCTM(context, 0, size);
    CGContextScaleCTM(context, 1, -1);
    CGContextSetShouldAntialias(context, true);

    CGMutablePathRef tile = CGPathCreateMutable();
    CGPathMoveToPoint(tile, NULL, 216, 38);
    CGPathAddLineToPoint(tile, NULL, 808, 38);
    CGPathAddCurveToPoint(tile, NULL, 906, 38, 986, 118, 986, 216);
    CGPathAddLineToPoint(tile, NULL, 986, 808);
    CGPathAddCurveToPoint(tile, NULL, 986, 906, 906, 986, 808, 986);
    CGPathAddLineToPoint(tile, NULL, 216, 986);
    CGPathAddCurveToPoint(tile, NULL, 118, 986, 38, 906, 38, 808);
    CGPathAddLineToPoint(tile, NULL, 38, 216);
    CGPathAddCurveToPoint(tile, NULL, 38, 118, 118, 38, 216, 38);
    CGPathCloseSubpath(tile);
    CGContextSaveGState(context);
    CGContextAddPath(context, tile);
    CGContextClip(context);
    CGColorRef blueLight = color(29, 103, 252, 255);
    CGColorRef blueDark = color(5, 45, 181, 255);
    const CGFloat stops[] = {0, 1};
    CFArrayRef colors = CFArrayCreate(NULL, (const void *[]){blueLight, blueDark}, 2, &kCFTypeArrayCallBacks);
    CGGradientRef gradient = CGGradientCreateWithColors(space, colors, stops);
    CGContextDrawLinearGradient(context, gradient, CGPointMake(230, 40), CGPointMake(790, 990), 0);
    CGContextRestoreGState(context);

    const CGFloat x = 502, y = 525;
    CGColorRef cream = color(255, 249, 226, 255);
    CGColorRef shadow = color(0, 20, 98, 95);
    CGContextSaveGState(context);
    CGContextSetShadowWithColor(context, CGSizeMake(9, 17), 23, shadow);
    CGContextSetFillColorWithColor(context, cream);
    ellipse(context, x, y, 340, true);
    CGContextRestoreGState(context);

    CGColorRef grooveDark = color(212, 197, 157, 200);
    CGColorRef grooveLight = color(255, 255, 241, 240);
    CGFloat radii[] = {303, 264, 225, 186};
    for (int i = 0; i < 4; i++) {
        CGContextSetStrokeColorWithColor(context, grooveDark);
        CGContextSetLineWidth(context, 5);
        ellipse(context, x, y, radii[i], false);
        CGContextSetStrokeColorWithColor(context, grooveLight);
        CGContextSetLineWidth(context, 3);
        ellipse(context, x, y, radii[i] - 5, false);
    }

    CGColorRef labelBlue = color(6, 54, 190, 255);
    CGContextSetFillColorWithColor(context, labelBlue);
    ellipse(context, x, y, 111, true);
    CGContextSetFillColorWithColor(context, cream);
    ellipse(context, x, y, 26, true);

    CGColorRef cutBlue = color(8, 65, 213, 255);
    CGContextSetFillColorWithColor(context, cutBlue);
    CGContextBeginPath(context);
    CGContextMoveToPoint(context, 568, 434);
    CGContextAddLineToPoint(context, 740, 262);
    CGContextAddLineToPoint(context, 758, 278);
    CGContextAddLineToPoint(context, 581, 448);
    CGContextClosePath(context);
    CGContextFillPath(context);

    CGColorRef lime = color(204, 255, 42, 255);
    CGContextSetFillColorWithColor(context, lime);
    CGContextBeginPath(context);
    CGContextMoveToPoint(context, 790, 193);
    CGContextAddQuadCurveToPoint(context, 793, 224, 805, 234);
    CGContextAddQuadCurveToPoint(context, 820, 246, 846, 249);
    CGContextAddQuadCurveToPoint(context, 820, 252, 805, 264);
    CGContextAddQuadCurveToPoint(context, 793, 278, 790, 305);
    CGContextAddQuadCurveToPoint(context, 787, 278, 775, 264);
    CGContextAddQuadCurveToPoint(context, 760, 252, 734, 249);
    CGContextAddQuadCurveToPoint(context, 760, 246, 775, 234);
    CGContextAddQuadCurveToPoint(context, 787, 224, 790, 193);
    CGContextClosePath(context);
    CGContextFillPath(context);

    CGImageRef image = CGBitmapContextCreateImage(context);
    CFURLRef destinationURL = CFURLCreateFromFileSystemRepresentation(NULL, (const UInt8 *)argv[1], strlen(argv[1]), false);
    CGImageDestinationRef destination = CGImageDestinationCreateWithURL(destinationURL, CFSTR("public.png"), 1, NULL);
    CGImageDestinationAddImage(destination, image, NULL);
    bool saved = CGImageDestinationFinalize(destination);

    CFRelease(destination);
    CFRelease(destinationURL);
    CGImageRelease(image);
    CGColorRelease(blueLight); CGColorRelease(blueDark);
    CGColorRelease(cream); CGColorRelease(shadow);
    CGColorRelease(grooveDark); CGColorRelease(grooveLight);
    CGColorRelease(labelBlue); CGColorRelease(cutBlue); CGColorRelease(lime);
    CGGradientRelease(gradient); CFRelease(colors); CGPathRelease(tile);
    CGContextRelease(context); CGColorSpaceRelease(space);
    return saved ? 0 : 1;
}
