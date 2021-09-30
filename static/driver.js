
function goEdit() {
    const urlParams = new URLSearchParams(window.location.search);
    const server = urlParams.get('server');
    window.location.href="/edit/"+server
}
