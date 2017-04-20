# Simple Session

A simple x11 session, just open a session, incompatible with xdg standard.


## DBus Interface

**Dest**: `org.jouyouyun.SimpleSession`

**Path**: `/org/jouyouyun/SimpleSession`

**Interface**: `org.jouyouyun.SimpleSession`

### Method

- Logout()
> Exit session

- Launch(cmd string)
> Launch applications, such as: nautilus --no-desktop

- ToggleDebug()
> Join debug mode

- SetOutputs(outputs string)
> Set outputs properties, such as: '{"outputs":[{"Output": "eDP1", "X":0, "Y": 0, "Width": 1920, "Height":1080}]}'


## Config

### Basic Config

You can custom the window manager, autostart scripts directory and background file.

If no config, nothing to do.

Example:
``` json
{
    "WM": "openbox",
    "autoscripts": "/etc/simple-session/autoscripts/",
    "background": "/usr/share/simple-session/background.jpg"
}
```


### Display Config

You can custom the primary priority, display mode and output blacklist, the output in blacklist will be off.
The outputs will be sorted by priority.

Display mode available values: extend, mirror, default: extend.

Example:
``` json
{
    "Priority": ["eDP1","LVDS1"],
    "Blacklist": ["VGA-2","DP2"],
    "Mode": "extend"
}
```


### Keybinding Config

You can custom shortcut list, must contain shortcut and action.

Example:
``` json
{
    "List":[
        {
            "Shortcut": "Super-T",
            "Action":"xterm"
        },
        {
            "Shortcut": "Super-Delete",
            "Action":"dbus-send --print-reply --dest=org.jouyouyun.SimpleSession /org/jouyouyun/SimpleSession org.jouyouyun.SimpleSession.Logout"
        }
    ]
}
```
