(function () {
    var app = new Vue({
        el: '#app',
        data: {
            dragging: false,
            images: [],
            nowImage: null
        },
        mounted: function () {
            var self = this;
            fetch('/lgtm').then(function (resp) {
                return resp.json();
            }).then(function (json) {
                self.images = json.images;
                self.$el.style.display = 'block';
            });
        },
        methods: {
            onSelectFile: function (e) {
                this.uploadFile(e.srcElement.files[0]);
            },
            onDragFile: function (e) {
                e.preventDefault();
                e.stopPropagation();
                this.dragging = true;
            },
            onDragLeaveFile: function (e) {
                e.preventDefault();
                e.stopPropagation();
                this.dragging = false;
            },
            onDropFile: function (e) {
                e.preventDefault();
                var file = e.dataTransfer.files[0];
                this.uploadFile(file);
            },
            uploadFile: function (file) {
                var self = this;
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
                this.nowImage = image;
            },
            onCloseModal: function (e) {
                e.preventDefault();
                this.nowImage = null;
            }
        }
    });
})();

