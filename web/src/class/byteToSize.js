// 字节单位转换
const bytesToSize = (bytes) => {
    if (bytes === 0) return '0 B';
    var k = 1024,
        sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
        i = Math.floor(Math.log(bytes) / Math.log(k));

    return Number((bytes / Math.pow(k, i)).toPrecision(3)).toLocaleString() + ' ' + sizes[i];
}
export default bytesToSize;