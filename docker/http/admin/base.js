var needleApp = {
    RESPONSE_OK: 'ok',
    RESPONSE_ERROR: 'error',
    _writer: null,
    setWriter: function (writer) {
        this._writer = writer;
    },
    /**
     * @return writer
     */
    getWriter: function () {
        if (this._writer == null) {
            throw new Error("No writer injected!");
        }
        return this._writer;
    },
    /**
     * Load list
     * @param pageNr
     */
    loadList: function (pageNr) {
        var perPage = 20;
        pageNr = (pageNr === undefined) ? 1 : pageNr;
        var offset = (pageNr * perPage) - perPage;
        $.ajax({
            url: "/url?offset=" + offset,
            success: function (response) {
                needleApp.getWriter().writeListResponse(response, perPage, offset, pageNr);
            }
        });
    },
    addNew: function (fromUrl, toUrl) {
        $.ajax({
            url: "/url",
            method: "POST",
            data: JSON.stringify({to_url: toUrl, from_url: fromUrl}),
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function (response) {
                needleApp.getWriter().writeNewResponse(response);
            },
            error: function (xhr, messageStatus, thrownError) {
                needleApp.getWriter().writeNewResponse(xhr.responseJSON);
            }
        });
    },
    updateOne: function (fromUrl, toUrl) {
        $.ajax({
            url: "/url",
            method: "PATCH",
            data: JSON.stringify({to_url: toUrl, from_url: fromUrl}),
            contentType: "application/json; charset=utf-8",
            dataType: "json",
            success: function (response) {
                needleApp.getWriter().writeUpdateResponse(response);
            },
            error: function (xhr, messageStatus, thrownError) {
                needleApp.getWriter().writeUpdateResponse(xhr.responseJSON);
            }
        });
    },
    removeOne: function (fromUrl) {
        if (!confirm("Are you sure you want to delete " + fromUrl + "?")) {
            return;
        }
        $.ajax({
            url: "/url/" + fromUrl,
            method: "DELETE",
            success: function (response) {
                needleApp.getWriter().writeRemoveResponse(response, fromUrl);
            },
            error: function (xhr, messageStatus, thrownError) {
                needleApp.getWriter().writeRemoveResponse(xhr.responseJSON, fromUrl);
            }
        });
    }
};