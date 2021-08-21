
var colorFunc = function (row, type, set, meta) {
    var index = meta.col - 12;
    var value = row["color" + index];

    var c = null;
    if (value < 0) {
        c = dyes[value+55]
    } else {
        c = colors[value];
    }

    if (c) {
        return "<div class='swatch' style='background-color: "
            + c.color + ";'></div> " + c.id + ": " + c.name;
    } else {
        if (value != 0) {
            console.log("Unknown color value: " + value + " for slot " + index);
        }
        return "<div class='swatch'></div> " + value + ": Unused";
    }
}

var showStats = function () {
    var showBase = $( '#showBase' )[0].checked;
    var showTamed = $( '#showTamed' )[0].checked;
    var showTotal = $( '#showTotal' )[0].checked;
    var showCurrent = $( '#showCurrent' )[0].checked;
    var showColor = $( '#showColor' )[0].checked;

    table.columns(3).visible(showBase);
    table.columns(4).visible(showTamed);
    table.columns(5).visible(showTotal);

    table.columns([12, 13, 14, 15, 16, 17]).visible(showColor);
    table.columns([18, 19, 20, 21, 22, 23, 24, 25]).visible(showCurrent);
    table.columns([26, 27, 28, 29, 30, 31, 32, 33]).visible(showBase);
    table.columns([34, 35, 36, 37, 38, 39, 40, 41]).visible(showTamed);
    table.columns([42, 43, 44, 45, 46, 47, 48, 49]).visible(showTotal);
};

var columns = [
    {"data": "name", "title": "Name"},
    {"data": "world", "title": "World", "visible": false},
    {"data": "class_name", "title": "Class"},
    {"data": "levels_wild", "title": "Base Lvl", "visible": false},
    {"data": "levels_tamed", "title": "Tame Lvl", "visible": false},
    {"data": "levels_total", "title": "Total Lvl"},

    {"data": "is_cryo", "title": "Stored?", "visible": false},
    {"data": "parent_class", "title": "Container Type", "visible": false},
    {"data": "parent_name", "title": "Container Name", "visible": false},

    {"data": "x", "visible": false},
    {"data": "y", "visible": false},
    {"data": "z", "visible": false},

    {"data": colorFunc, "title": "C0", "searchBuilderTitle": "Color 0", "visible": false},
    {"data": colorFunc, "title": "C1", "searchBuilderTitle": "Color 1", "visible": false},
    {"data": colorFunc, "title": "C2", "searchBuilderTitle": "Color 2", "visible": false},
    {"data": colorFunc, "title": "C3", "searchBuilderTitle": "Color 3", "visible": false},
    {"data": colorFunc, "title": "C4", "searchBuilderTitle": "Color 4", "visible": false},
    {"data": colorFunc, "title": "C5", "searchBuilderTitle": "Color 5", "visible": false},

    {"data": "health_current", "title": "H", "searchBuilderTitle": "Current Health", "visible": false},
    {"data": "stamina_current", "title": "St", "searchBuilderTitle": "Current Stamina", "visible": false},
    {"data": "torpor_current", "title": "T", "searchBuilderTitle": "Current Torpor", "visible": false},
    {"data": "oxygen_current", "title": "O", "searchBuilderTitle": "Current Oxygen", "visible": false},
    {"data": "food_current", "title": "F", "searchBuilderTitle": "Current Food", "visible": false},
    {"data": "weight_current", "title": "W", "searchBuilderTitle": "Current Weight", "visible": false},
    {"data": "melee_current", "title": "M", "searchBuilderTitle": "Current Melee", "visible": false},
    {"data": "speed_current", "title": "Sp", "searchBuilderTitle": "Current Speed", "visible": false},

    {"data": "health_wild", "title": "H", "searchBuilderTitle": "Base Health", "visible": false},
    {"data": "stamina_wild", "title": "St", "searchBuilderTitle": "Base Stamina", "visible": false},
    {"data": "torpor_wild", "title": "T", "searchBuilderTitle": "Base Torpor", "visible": false},
    {"data": "oxygen_wild", "title": "O", "searchBuilderTitle": "Base Oxygen", "visible": false},
    {"data": "food_wild", "title": "F", "searchBuilderTitle": "Base Food", "visible": false},
    {"data": "weight_wild", "title": "W", "searchBuilderTitle": "Base Weight", "visible": false},
    {"data": "melee_wild", "title": "M", "searchBuilderTitle": "Base Melee", "visible": false},
    {"data": "speed_wild", "title": "Sp", "searchBuilderTitle": "Base Speed", "visible": false},

    {"data": "health_tamed", "title": "H", "searchBuilderTitle": "Health Tamed Points", "visible": false},
    {"data": "stamina_tamed", "title": "St", "searchBuilderTitle": "Stamina Tamed Points", "visible": false},
    {"data": "torpor_tamed", "title": "T", "searchBuilderTitle": "Torpor Tamed Points", "visible": false},
    {"data": "oxygen_tamed", "title": "O", "searchBuilderTitle": "Oxygen Tamed Points", "visible": false},
    {"data": "food_tamed", "title": "F", "searchBuilderTitle": "Food Tamed Points", "visible": false},
    {"data": "weight_tamed", "title": "W", "searchBuilderTitle": "Weight Tamed Points", "visible": false},
    {"data": "melee_tamed", "title": "M", "searchBuilderTitle": "Melee Tamed Points", "visible": false},
    {"data": "speed_tamed", "title": "Sp", "searchBuilderTitle": "Speed Tamed Points", "visible": false},

    {"data": "health_total", "title":"H", "searchBuilderTitle": "Total Health"},
    {"data": "stamina_total", "title":"St", "searchBuilderTitle": "Total Stamina"},
    {"data": "torpor_total", "title":"T", "searchBuilderTitle": "Total Torpor"},
    {"data": "oxygen_total", "title":"O", "searchBuilderTitle": "Total Oxygen"},
    {"data": "food_total", "title":"F", "searchBuilderTitle": "Total Food"},
    {"data": "weight_total", "title":"W", "searchBuilderTitle": "Total Weight"},
    {"data": "melee_total", "title":"M", "searchBuilderTitle": "Total Melee"},
    {"data": "speed_total", "title":"Sp", "searchBuilderTitle": "Total Speed"}
]

var tableOptions = {
    "ajax": {"url":"api/dinos", "dataSrc":""},
    "columns": columns,

    "dom": "rtpil",
    "pageLength": 50,
    "scrollX": true,

    "language": {
        "searchBuilder": {
            "title": ""
        },
    },

    "select": {
        "info": false,
        "style": "single",
    },

    "searchBuilder": {
        "columns": [0, 1, 2, 3, 4, 5, 6, 7, 8, 12, 13, 14, 15, 16, 17,
                    18, 19, 20, 21, 22, 23, 24, 25,
                    26, 27, 28, 29, 30, 31, 32, 33,
                    34, 35, 36, 37, 38, 39, 40, 41,
                    42, 43, 44, 45, 46, 47, 48, 49]
    }
}
