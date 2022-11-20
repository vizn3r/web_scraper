# Simple Web Scraper

## How to use

### Windows

1. Download the `web_scrapper.exe` file
2. Run `web_srapper.exe`
3. Enter `web url`
4. Enter `element`

    - This could be `html` tag
    - This could be `class`
    - This could be `id`
    - Or everything combined

    - If you want to search for more elements, separate them by `space`
    - Example:

        ```txt
        // single element examples
        Enter 'element/s': #wrapper2 --> id
        Enter 'element/s': .Corner14 --> class
        Enter 'element/s': img --> html tag

        // multi element examples
        Enter 'element/s': #wrapper2 .Corner14 --> look for elements with 'wrapper2' id or 'Corner14' class
        Enter 'element/s': img .Corner16 --> look for elements with 'img' tag or 'Corner16' class

        // specific element with multiple classes / id's / elements
        Enter 'elemen/s': .BreakInsideAvoid.frameLight.Corner14 --> look for elements with classes 'BreakInsideAvoid', 'frameLight' and 'Corner14'
        ```

5. Enter `output directory`

    - If you want the same directory that the executable is, press `enter`
    - If you want to create new directory, you need to enter `relative path` to the executable
    - ! YOU NEED TO END THE DIRECTORY WITH `/` !
    - Example:

        ```txt
        // same dir as the executable
        Enter 'output directory': <press enter>

        // relative dir to executable
        Enter 'output directory': out/ --> will create 'out' dir in directory that the executable is
        Enter 'output directory': ../out/ --> will create 'out' dir in directory before that directory the executable is
        ```

- If you have any problems/questions, create New Issue
