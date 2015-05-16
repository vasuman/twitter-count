function callAPI(method, params, success, fail) {
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '/api/' + method);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.onload = function() {
                if (this.status != 200) {
                        console.log("API call error!");
                        fail(this.responseText);
                        return;
                }
                var resp = JSON.parse(this.responseText);
                success(resp);
        };
        xhr.send(JSON.stringify(params));
}
