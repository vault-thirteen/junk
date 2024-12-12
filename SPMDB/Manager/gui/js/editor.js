// Settings.
const backendHost = "localhost";
const backendPort = 2000;
const backendAddress = "http://" + backendHost + ":" + backendPort;

// URLs of handlers.
const huGetMoviesCount = "/movies/count";
const huGetMovie = "/movie";
const huGetServer = "/server";

// Ids of elements.
videoHolderId = "videoHolder";
videoContentId = "videoContent";

// Objects.
var videoHolder;
var videoContent;

// JavaScript ShitBox.
var moviesCount;
var movieById;
var serverById;

function initPage() {
    videoHolder = document.getElementById(videoHolderId);
    videoContent = document.getElementById(videoContentId);
    videoContent.addEventListener("loadeddata", videoIsLoaded, true);
    videoContent.addEventListener("ended", videoHasEnded, true);

    window.onresize = scaleVideo;

    getMoviesCount(); // moviesCount <- x.
}

function videoIsLoaded() {
    scaleVideo();
    //videoContent.play();
}

function videoHasEnded() {
}

function getMoviesCount() {
    let nextFn = initPage_part2;
    let result;
    let xhr = new XMLHttpRequest();
    let url = backendAddress + huGetMoviesCount;
    xhr.open('GET', url);

    xhr.onload = function () {
        if (xhr.status !== 200) {
            result = 0;
            console.error(xhr.status + " " + url)
        } else {
            result = parseInt(xhr.responseText);
        }

        moviesCount = result;
        nextFn();
    };

    xhr.send();
}

function initPage_part2() {
    let rndMovieId = getRandomMovieId();
    getMovieById(rndMovieId); // movieById <- x.
}

function getRandomMovieId() {
    return 1 + Math.floor(Math.random() * moviesCount);
}

function getMovieById(id) {
    let nextFn = initPage_part3;
    let result;
    let xhr = new XMLHttpRequest();
    let url = backendAddress + huGetMovie + "?id=" + id;
    xhr.open('GET', url);

    xhr.onload = function (e) {
        if (xhr.status !== 200) {
            console.error(xhr.status + " " + url)
        } else {
            result = JSON.parse(xhr.responseText);
        }

        movieById = result;
        nextFn();
    };

    xhr.send();
}

function initPage_part3() {
    console.log(movieById);
    let serverId = movieById.FileServerId;
    getServerById(serverId); // serverById <- x.
}

function getServerById(id) {
    let nextFn = initPage_part4;
    let result;
    let xhr = new XMLHttpRequest();
    let url = backendAddress + huGetServer + "?id=" + id;
    xhr.open('GET', url);

    xhr.onload = function (e) {
        if (xhr.status !== 200) {
            console.error(xhr.status + " " + url)
        } else {
            result = JSON.parse(xhr.responseText);
        }

        serverById = result;
        nextFn();
    };

    xhr.send();
}

function initPage_part4() {
    console.log(serverById);
    let url = serverById.Address + movieById.FilePath + "/" + movieById.FileName;
    let contentType = getMimeTypeForFileExtension(movieById.FileExtension);

    while (videoContent.firstChild) {
        videoContent.removeChild(videoContent.firstChild);
    }

    var source = document.createElement('source');
    source.setAttribute("src", url);
    source.setAttribute("type", contentType);

    videoContent.appendChild(source);
    videoContent.load();
}

function getScaleKoefficient(w, wMax, h, hMax) {
    var kW = wMax / w;
    var kH = hMax / h;
    var k;

    if (kH < kW) {
        k = kH;
    } else {
        k = kW;
    }

    return k;
}

function scaleVideo() {
    // Scale the Video.
    var page_height_visible = videoHolder.offsetHeight;
    var page_width_visible = videoHolder.offsetWidth;
    var video_source_height = videoContent.videoHeight;
    var video_source_width = videoContent.videoWidth;
    var scaleRatio = getScaleKoefficient(
        video_source_width,
        page_width_visible,
        video_source_height,
        page_height_visible
    );

    var widthNew = videoContent.videoWidth * scaleRatio;
    var heightNew = videoContent.videoHeight * scaleRatio;

    //console.log(widthNew + " x " + heightNew);

    videoContent.width = widthNew;
    videoContent.height = heightNew;
}

function getMimeTypeForFileExtension(ext) {
    switch (ext) {
        case "avi":
            return "video/x-msvideo";
        case "f4v":
            return "video/x-f4v";
        case "m2ts":
            return "video/mp2t"; // Works only in Safari.
        case "m4v":
            return "video/x-m4v";
        case "mkv":
            return "video/x-matroska";
        case "mov":
            return "video/quicktime";
        case "mp4":
            return "video/mp4";
        case "mpeg":
            return "video/mpeg";
        case "mpg":
            return "video/mpeg";
        case "ts":
            return "video/mp2t"; // Works only in Safari.
        case "wmv":
            return "video/x-ms-wmv";
    }
}
