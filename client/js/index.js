let viewer;

let errorCodes = {
    1: 'An unknown failure has occurred.',
    2: 'Bad data (corrupted or malformed) was encountered.',
    3: 'A network failure was encountered.',
    4: 'Access was denied to a network resource (HTTP 403)',
    5: 'A network resource could not be found (HTTP 404)',
    6: 'A server error was returned when accessing a network resource (HTTP 5xx)',
    7: 'An unhandled response code was returned when accessing a network resource (HTTP "everything else")',
    8: 'Browser error: webGL is not supported by the current browser',
    9: 'There is nothing viewable in the fetched document',
    10: 'Browser error: webGL is supported, but not enabled',
    11: 'There is no geomtry in loaded model',
    12: 'Collaboration server error'
};


fetch("/gettoken")
    .then(res => res.json())
    .then(result => {
        let token = result["access_token"];
        let options = {
            env: 'AutodeskProduction',
            accessToken: token
        };
        Autodesk.Viewing.Initializer(options, function onInitialized() {

            fetch("/geturn")
                .then(res => res.json())
                .then(result => {
                    let urn = result["urn"];
                    if ( urn.length === 0) {
                        window.location.replace('/upload.html');
                        return;
                    }
                    let documentId = 'urn:' + urn;
                    // documentId = "urn:dXJuOmFkc2sub2JqZWN0czpvcy5vYmplY3Q6YnVja2V0Mzk0MjMzOTQzNDU3ODE0OTYzNi9IVzIucnZ0";
                    Autodesk.Viewing.Document.load(documentId, onDocumentLoadSuccess, onDocumentLoadFailure);
                })
                .catch(error => console.log("Could not fetch URN: ",error))
        });
    })
    .catch(error => console.log("Could not fetch token: ", error));


/**
 * Autodesk.Viewing.Document.load() success callback.
 * Proceeds with model initialization.
 */
function onDocumentLoadSuccess(doc) {

    // A document contains references to 3D and 2D viewables.
    let viewables = Autodesk.Viewing.Document.getSubItemsWithProperties(doc.getRootItem(), {'type': 'geometry'}, true);
    if (viewables.length === 0) {
        console.error('Document contains no viewables.');
        return;
    }

    // Choose any of the avialble viewables
    let initialViewable = viewables[0];
    let svfUrl = doc.getViewablePath(initialViewable);
    let modelOptions = {
        sharedPropertyDbPath: doc.getPropertyDbPath()
    };

    let viewerDiv = document.getElementById('viewerDiv');
    viewer = new Autodesk.Viewing.Private.GuiViewer3D(viewerDiv);
    viewer.start(svfUrl, modelOptions, onLoadModelSuccess, onLoadModelError);
}

/**
 * Autodesk.Viewing.Document.load() failuire callback.
 */
function onDocumentLoadFailure(viewerErrorCode) {
    console.error('onDocumentLoadFailure() - errorCode:'
        + viewerErrorCode
        + " ==> "
        + errorCodes[viewerErrorCode]);
}

/**
 * viewer.loadModel() success callback.
 * Invoked after the model's SVF has been initially loaded.
 * It may trigger before any geometry has been downloaded and displayed on-screen.
 */
function onLoadModelSuccess(model) {
    console.log('onLoadModelSuccess()!');
    console.log('Validate model loaded: ' + (viewer.model === model));
    console.log(model);
}

/**
 * viewer.loadModel() failure callback.
 * Invoked when there's an error fetching the SVF file.
 */
function onLoadModelError(viewerErrorCode) {
    console.error('onLoadModelError() - errorCode:' + viewerErrorCode);
}