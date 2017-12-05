function Writer() {

    this._writeEmptyList = function () {
        $('#list').append($('<tr></tr>').html("<td colspan='4'>No url's defined!</td>"));
    };

    /**
     * Actually format elements and write them to list
     */
    this._factoryListElement = function (element, index) {
        var mainContainer = $('<tr></tr>').attr('id', element.from_url);
        // index
        mainContainer.append(
            $('<th></th>').attr('scope', 'row').html(
                parseInt(index) + 1
            )
        );
        // build link
        var BASE_URL = config.base_url;
        var link = BASE_URL + element["from_url"];

        mainContainer.append(
            $('<td></td>').html(
                "<a href='" + link + "' target='_blank'><span>" + BASE_URL + "</span>" + element.from_url + "</a>"
            )
        );
        mainContainer.append(
            $('<td></td>').html(element.to_url)
        );
        mainContainer.append(
            $('<td></td>').html(
                "<button class='btn btn-primary btn-sm' onclick=\"needleApp.removeOne('" + element.from_url + "')\">Delete</button>&nbsp" +
                "<button class='btn btn-info btn-sm' onclick=\"showEdit('" + element.from_url + "')\">Edit</button>"
            )
        );
        return mainContainer;
    };

    /**
     * Write pagination elements
     */
    this._writePagination = function (totalItems, perPage, pageNr) {
        var pageCount = Math.ceil(totalItems / perPage);
        $('#pagination').html('').append(
            $('<li></li>').addClass('page-item disabled').html('<a class="page-link" href="#">Jump to page:</a>')
        );

        for (var i = 1; i <= pageCount; i++) {
            var element = $('<li></li>').addClass('page-item').html(
                '<a class="page-link" href="#" onclick="needleApp.loadList(\'' + i + '\')">' + i + '</a>'
            );
            if (i == pageNr) {
                element.addClass('active');
            }
            $('#pagination').append(element);
        }
    };

    /**
     * Write reponse for listing
     * @param response
     * @param perPage
     * @param offset
     * @param pageNr
     */
    this.writeListResponse = function(response, perPage, offset, pageNr) {
        if (response["total"] == 0) {
            this._writeEmptyList();
            return;
        }

        $('#list').html('');
        for (var item in response["list"]) {
            $('#list').append(this._factoryListElement(response["list"][item], item));
        }
        // build pagination
        this._writePagination(response['total'], perPage, pageNr);
    };

    this.writeNewResponse = function (response) {
        if (response.status == needleApp.RESPONSE_ERROR) {
            $('#alert-modal').find('p').html(response.message);
            $('#alert-modal').modal('show');
            return;
        }
        $('#update_url').modal('hide');
        needleApp.loadList($('#page-item active a').html());
    };

    this.writeUpdateResponse = function (response) {
        if (response.status == needleApp.RESPONSE_ERROR) {
            $('#alert-modal').find('p').html(response.message);
            $('#alert-modal').modal('show');
            return;
        }

        $('#alert-modal').find('p').html("Update was successful!");
        $('#alert-modal').modal('show');
        $('#update_url').modal('hide');
        needleApp.loadList($('#page-item active a').html());
    };

    /**
     * Handle remove response
     * @param response
     * @param itemReferenceId
     */
    this.writeRemoveResponse = function (response, itemReferenceId) {
        if (response.status == needleApp.RESPONSE_ERROR) {
            $('#alert-modal').find('p').html(response.message);
            $('#alert-modal').modal('show');
            return;
        }
        $('#' + itemReferenceId).fadeOut(function () {
            $(this).remove();
        });
    };
    return this;
}

function addNewElement() {
    needleApp.addNew($('#from_url').val(), $('#to_url').val());
}

function updateElement() {
    needleApp.updateOne($('#from_url').val(), $('#to_url').val());
}

function showEdit(item) {
    $.ajax({
        url: "/url?item=" + item,
        success: function (response) {
            urlItem = response.list.pop();

            $('#from_url').val(urlItem.from_url);
            $('#to_url').val(urlItem.to_url);
            $('.rel-add').addClass('hidden');
            $('.rel-update').removeClass('hidden');

            $('#from_url').attr('disabled', 'disabled');
            $('#update_url').modal('show');
        }
    });
}

function showAddModal() {
    $('.rel-add').removeClass('hidden');
    $('.rel-update').addClass('hidden');
    $('#update_url').modal('show');
    $('#from_url').removeAttr('disabled');
}

$(document).ready(function() {
    var writer = new Writer();
    needleApp.setWriter(writer);
    needleApp.loadList();
});