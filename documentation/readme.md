# TELKO-MOMENT-PROJECT 🗄

## Bio 💁

>  we're just a couple of college kids that wanted to get a hold of flutter and also get to use nodejs and some multi-service architectures  on something that ***'works-in-production'***.
>  So we just rallied ourselves and then started a project (_this project_) so that wwe could boost our skills. We also wish to use a few if not all (_if possible_) good software engineering principles to have a good understanding of how they work, experience then also and see them in action.

## Branch-strategy 🌳

> Since the project is quite small and we also have a small team, we found suitable to have to use the **GitHub Flow Branch Strategy**.
> - Putting it in simple words, the **github flow branch strategy** has a ***main-branch*** that is branched off by a ***feature-branch*** whenever there is a new feature being implemented.
> - The ***feature-branch*** will then be merged back with the ***main-branch*** whenever the said feature has been completed, checked and ready for implementation.
> - but as for us we will be tweaking the branching strategy to have an **Extended GitHub Flow Branch Strategy**.
>   +   the extended version (_our version_) will have an additional branch between the ***main-branch*** and ***feature-branch***, and this will be the ***develop-branch***.
>   +   the ***develop-branch*** will be home to the current in-development code as most code will be pushed from the feature branches.
>   +   this will keep the main branch clean and also reduce errors in production code.
> - The branching strategy is quite nice and simple :) .


## Description 🛈

> - server and mobile-app for a chat-app.
> - it depends on sockets _(over stomp protocol)_ and rest-API for serving messages and files.
> - inspired by other open-source chat-apps that we use day to day.
> - the base architectures will be described in the root folder and then for specific architectures it would have to be added to the specific folders.
> - this will be the **root** folder for the blueprint and documentation of the project and the 2 apps
> - **SERVER-APP** will be made readily available on <https://github.com/Stroustrups-Sentinel/telko-moment-server>
> - **MOBILE-APP** will be kept with us for now :D.



## Databases 💾

> - we mostly wanted to have a grasp and also lie in the _'safety net'_ of a NoSql database since you could escape _'scar free'_ when you all of a sudden change a schema and also even in production as there would be lesser overhead as compared to a SQL database.
> Besides the safety net, we also wanted to have a better understanding of the NoSql database(s) by using them in a fairly _'larger'_ app compared to our apps that we mostly develop in school due to the smaller deadlines, lazybones in teams 😅 and other reasons to make me convince you that the project will be small.

## Softwares 💻

> These are the softwares that we will be using in this project.
> Most of them will be for adding productivity and others for improving _quality of life_ as we work on the project, although the plugins are not compulsory just in case you have an alternative in mind.

> The list is as follows :
 1. VisualStudio-Code 
    -   extensions :
          +   Material Icon Theme
          +   Beautify
          +   Prettier
          +   Better Comments
          +   Bracket Pair Colorizer (Although it was added internally on new versions)
          +   Color Picker
          +   Code Runner
          +   DotENV
          +   Docker
          +   Code Spell Checker
          +   Flutter, Dart
          +   Gitlens
          +   *GitHub Theme (If you like the themes)
          +   Live Server 
          +   Markdown All in One
          +   Rainbow CSV
          +   Todo Tree
          +   XML
          +   YAML
          +   Path intellisense
          +   Npm Intellisense


 2. Android Studio year:2021+
    > 1. Command line tools
    > 2. Flutter ver:3 SDK (contains Dart SDK)
    >    * Run `cmd flutter doctor ` (check missing modules, consider the important ones like - command-line tools & android-licenses)
        
    > 3. Android Virtual Device: 
    >    * android-version: **5.0**
    >    * screen-size: **Pixel 4XL**
    >    * cpu-Architecture: **x86_x64**
    >    * api-level : **21**
    >    * google-play-store : **yes**
 
 3. Ghostwriter (QOL on Markdown, although there will be vscode markdown extension, more like the real-Radio vs Radio-mobile-app scenario )
 
 4. Github Desktop / GitKraken (git-gui app, choose the one that works best for you )
 
 5. Gaphor (UMLs)
 
 6. Figma (UMLs, User stories)
 
 7. Mongodb & Compass (for GUI)