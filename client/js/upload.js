+ function($) {
    'use strict';

    // UPLOAD CLASS DEFINITION
    // ======================

    var dropZone = document.getElementById('drop-zone');

    var startUpload = function(files) {
        var status = document.getElementById("status");
        fetch('/upload', {
           method: 'POST',
            headers: new Headers({
                "FileName": files[0].name
            }),
           body: files[0]
        }).then(function(response) {
            console.log(response.json());
            window.location.replace('/');
        });


        status.innerText = "Processing ..."
    };

    dropZone.ondrop = function(e) {
        e.preventDefault();
        this.className = 'upload-drop-zone';

        startUpload(e.dataTransfer.files)
    };

    dropZone.ondragover = function() {
        this.className = 'upload-drop-zone drop';
        return false;
    };

    dropZone.ondragleave = function() {
        this.className = 'upload-drop-zone';
        return false;
    }

}(jQuery);