const timeTransform = (second) => {
    var duration
    var days = Math.floor(second / 86400);
    var hours = Math.floor((second % 86400) / 3600);
    var minutes = Math.floor(((second % 86400) % 3600) / 60);
    var seconds = Math.floor(((second % 86400) % 3600) % 60);
    if (days > 0) duration = days + "天";
    else if (hours > 0) duration = hours + "小时";
    else if (minutes > 0) duration = minutes + "分钟";
    else if (seconds > 0) duration = seconds + "秒";
    return duration;
}

export default timeTransform