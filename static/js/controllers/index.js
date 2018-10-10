'use strict';

function initFun () {
  var container = document.getElementById("doc-detail-container")
  var status = document.getElementById("poe-status")
  var hash = document.getElementById("hash-info")
  if (status == null) {
    container.style.display = "none"
    hash.style.display = "block"
  } else {
    container.style.display = "block"
    hash.style.display = "none"
  }
}

var isInteger = function(n) {
    return n == parseInt(n, 10);
}

$(document).ready(function() {

  var bar = $('.bar');
  var upload_form = $('#upload_form');
  var explain = $('#explain');
  var dropbox = $('.dropbox');
  var selectedFile;
  var html5 = window.File && window.FileReader && window.FileList && window.Blob;
  $('#wait').hide();

  var handleFileSelect = function(f) {
    if (!html5) {
      return;
    }
    explain.html(translate('正在加载...'));
    var output = '';
    output = translate('Preparing to hash ') + f.name + ' (' + (f.type || translate('n/a')) + ') - '
      + f.size + translate(' bytes, ') + translate('last modified: ')
      + (f.lastModifiedDate ? f.lastModifiedDate
      .toLocaleDateString() : translate('n/a')) + '';

    var reader = new FileReader();
    reader.onload = function(e) {
      var data = e.target.result;
      bar.width(0 + '%');
      bar.addClass('bar-success');
      explain.html(translate('正在生成... ') + translate('Initializing'));
      setTimeout(function() {
        selectedFile = f
        CryptoJS.SHA256(data, crypto_callback, crypto_finish);
      }, 200);

    };
    reader.onprogress = function(evt) {
      if (evt.lengthComputable) {
        var w = (((evt.loaded / evt.total) * 100).toFixed(2));
        bar.width(w + '%');
      }
    };
    reader.readAsBinaryString(f);
  };
  if (!html5) {
    explain.html(translate('disclaimer'));
    upload_form.show();
  } else {
    dropbox.show();
    dropbox.filedrop({
      callback: handleFileSelect
    });
    dropbox.click(function() {
      $('#poe_file').click();
    });
  }

  var crypto_callback = function(p) {
    var w = ((p * 100).toFixed(0));
    bar.width(w + '%');
    explain.html(translate('正在生成... ') + (w) + '%');
  };

  var crypto_finish = function(hash) {
    bar.width(100 + '%');
    explain.html(translate('Document hash: ') + hash);

    var docHash = (document.getElementById("explain").innerHTML).substring(15)
    var formData = new FormData(document.getElementById('upload_form'));
　　// 建立一个upload表单项，值为上传的文件
　　formData.append('poe_file', document.getElementById('poe_file').files[0]);
	 formData.append('hash', docHash);
	 $.ajax({
    async: false,
 		url: "http://localhost:8081/v1/upload",
 		type: "POST",
		contentType : false,
 		data: formData,
		dataType: "html",
		processData: false,
		headers: {
		 "Accept": "application/json",
		 "API-Key":"xxxx",
		 "Bc-Invoke-Mode": "sync"
   },
   success: function (data) {
     $("#allBody").html(data);//刷新整个body页面的html
   }
 	})
};

  document.getElementById('poe_file').addEventListener('change', function(evt) {
    var f = evt.target.files[0];
    handleFileSelect(f);
  }, false);

});
