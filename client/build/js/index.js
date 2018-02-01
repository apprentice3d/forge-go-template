var viewerApp;
var viewer = null;
var tree = null;
var options = {
    env: 'AutodeskProduction',
    getAccessToken: function (onGetAccessToken) {
        var token_fetcher = '/gettoken';
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.open("GET", token_fetcher, false);
        xmlHttp.send(null);
        var data = JSON.parse(xmlHttp.responseText);
        var accessToken = data["access_token"];
        var expireTimeSeconds = data["expires_in"];


        console.log("received token: " + accessToken);

        onGetAccessToken(accessToken, expireTimeSeconds);
    }

};

function getRecentURN() {
    var urn_fetcher = '/geturn';
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.open("GET", urn_fetcher, false);
    xmlHttp.send(null);
    var data = JSON.parse(xmlHttp.responseText);
    var urn = data["urn"]
    console.log("Received from server the following urn: " + urn + "of length " + urn.length);
    if ( urn.length === 0) {window.location.replace('/upload.html');}
    return "urn:" + urn;
}


var documentId = getRecentURN();

Autodesk.Viewing.Initializer(options, function onInitialized() {
    console.log("Using the following URN: " + documentId);
    viewerApp = new Autodesk.Viewing.ViewingApplication('viewerDiv');
    viewerApp.registerViewer(viewerApp.k3D, Autodesk.Viewing.Private.GuiViewer3D);
    viewerApp.loadDocument(documentId, onDocumentLoadSuccess, onDocumentLoadFailure);
    viewer = viewerApp.getCurrentViewer();
});

function onDocumentLoadSuccess(doc) {

    // We could still make use of Document.getSubItemsWithProperties()
    // However, when using a ViewingApplication, we have access to the **bubble** attribute,
    // which references the root node of a graph that wraps each object from the Manifest JSON.
    var viewables = viewerApp.bubble.search({ 'type': 'geometry' });
    if (viewables.length === 0) {
        console.error('Document contains no viewables.');
        return;
    }

    // Choose any of the avialble viewables
    viewerApp.selectItem(viewables[0].data, onItemLoadSuccess, onItemLoadFail);


}

function onDocumentLoadFailure(viewerErrorCode) {
    console.log(documentId);
    console.error('onDocumentLoadFailure() - errorCode:' + viewerErrorCode);
}

function onItemLoadSuccess(reported_viewer, item) {

    viewer = reported_viewer;
    viewer.addEventListener(Autodesk.Viewing.OBJECT_TREE_CREATED_EVENT, setupMyModel);
}

function onItemLoadFail(errorCode) {
    console.error('onItemLoadFail() - errorCode:' + errorCode);
}



function setupMyModel() {
    /*============================ CUSTOMIZE MODEL HERE =======================*/

    tree = viewer.model.getData().instanceTree;


    /*============================ END OF CUSTOMIZATION =======================*/

}