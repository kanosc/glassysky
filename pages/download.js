$(document).ready(function(){
	$(".delete-btn").click(function(){
		var filename = $(this).attr("name");
		if (confirm("是否确认删除文件 " + filename) == false) {
			return;
		}
		var uname = $(this).attr("uname");
		var url = "/delete?filename=" + filename;
		$.ajax({url:url,
			error:function(xhr){
				alert("Error: " + xhr.status + " " + xhr.statusText + ", " + xhr.responseText)
			},
			success:function(res){
				$("#"+uname).remove();
			}
		});
	});
	$('#upload-form').submit(function(event) {
		event.preventDefault(); // 阻止表单默认提交行为

		var files = $('#file-input')[0].files;
		if (files.length == 0) {
			alert('请选择要上传的文件');
			return false;
		}

		// 创建XMLHttpRequest对象
		var xhr = new XMLHttpRequest();

		// 监听上传进度事件
		xhr.upload.addEventListener('progress', function(event) {
			if (event.lengthComputable) {
				var percent = Math.round((event.loaded / event.total) * 100);
				$('.progress-bar').css('width', percent + '%').attr('aria-valuenow', percent);
				$('#progress-text').text(percent + '%');
			}
		});

		// 监听上传完成事件
		xhr.addEventListener('load', function() {
			alert('上传成功，请刷新页面！');
			$('#progress-area').hide();
			//var responseHTML = xhr.responseText;
			//document.write(responseHTML);
			window.location.href="/download_list";

		});

		// 监听上传错误事件
		xhr.addEventListener('error', function() {
			alert('上传失败！');
			$('#progress-area').hide();
			var responseError = xhr.responseText;
			alert(responseError)

		});

		// 开始上传
		xhr.open('POST', '/upload');
		var formData = new FormData();
		for (var i = 0; i < files.length; i++) {
			formData.append('upload', files[i]);
		}
		xhr.send(formData);

		// 显示上传进度条
		$('#progress-area').show();
	});

	var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
	var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
		return new bootstrap.Tooltip(tooltipTriggerEl)
	});
});
