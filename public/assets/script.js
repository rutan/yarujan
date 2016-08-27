(function () {
    var app = new Vue({
        el: '#app',
        data: {
            dragging: false,
            images: []
        },
        methods: {
            onDragFile: function (e) {
                e.preventDefault();
                e.stopPropagation();
                this.dragging = true;
                console.log(e);
            },
            onDropFile: function (e) {
                var self = this;
                e.preventDefault();
                var file = e.dataTransfer.files[0];
                var formData = new FormData();
                formData.append("file", file);
                fetch('/lgtm', {
                    method: 'POST',
                    body: formData
                }).then(function (response) {
                    console.log(response);
                    return response.json();
                }).then(function (json) {
                    console.log(json);
                    self.images.unshift(json);
                    self.dragging = false;
                }).catch(function (e) {
                    console.error(e);
                    alert('failed generate LGTM');
                    self.dragging = false;
                });
            },
            onClickPicture: function (e, image) {
                e.preventDefault();
            }
        }
    });
})();

