
const addZero = (m) => {
    return m < 10 ? '0' + m : m;
}

const timestampToDate = (timestamp) => {
    var time = new Date(timestamp * 1000);
    var y = time.getFullYear();
    var M = time.getMonth() + 1;
    var d = time.getDate();
    var h = time.getHours();
    var m = time.getMinutes();
    var s = time.getSeconds();
    return y + '-' + addZero(M) + '-' + addZero(d) + ' ' + addZero(h) + ':' + addZero(m) + ':' + addZero(s);
}

export default timestampToDate