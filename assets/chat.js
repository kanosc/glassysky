var flagE = false;
$(document).ready(function(){

        /*
        var name = "";
        var name = prompt("please input your name for chat");
        if (name == "") {
                name = "Guest" + Math.floor(Math.random() * 1000);
        }*/

        var lockReconnect = false;
        //var url = "ws://" + window.location.host + "/wschat";
        var roomname = $("#roomname").text();
        var url = socketPrefix + "://" + window.location.host + "/wschat/" + roomname;
        var ws;
        createWebSocket();
        setInterval(createWebSocket, 30000)
        //var name = "Guest" + Math.floor(Math.random() * 1000);



        var now = function () {
                var date = new Date(new Date().getTime()+(parseInt(new Date().getTimezoneOffset()/60) + 8)*3600*1000).toString();
                var dateArr = date.split(" ");
                return dateArr[1] + " " + dateArr[2] + " " + dateArr[3] + " " + dateArr[4];
        };

        function init() {
                console.log("init begin");
                var name = $("#username").text();
                var roomname = $("#roomname").text();
                var chat = document.getElementById("chat");
                var text = document.getElementById("text");

                //var ws = new WebSocket(url);

                ws.onclose = function (e) {
                        flagE = false;
                        console.log('链接关闭' + e.code + new Date().getTime());
                        // 服务器关闭  不重连
                        if (e.code != 1006){
                        //      reconnect();
                                createWebSocket();
                        }
                };

                ws.onmessage = function (msg) {
                        var line =  msg.data;
                        //chat.innerText = line + chat.innerText;
                        function insertParagraph() {
                                // 创建一个新的<p>元素
                                const paragraph = document.createElement('div');
                                paragraph.textContent = line; // 设置<p>标签的内容为'hello'
                                paragraph.classList.add('alert');
                                paragraph.classList.add('alert-success');
                                paragraph.setAttribute("role", "alert");

                                // 在id为'chat'的元素中最前面插入<p>标签
                                const chatElement = document.getElementById('chat');
                                if (chatElement.firstChild) {
                                        chatElement.insertBefore(paragraph, chatElement.firstChild);
                                } else {
                                        chatElement.appendChild(paragraph);
                                }
                        }

                        // 调用函数来执行插入<p>标签的操作
                        insertParagraph();
                        chat.scrollTop = 0;
                };

                //text.onkeydown = sendMsg(e);
                text.onkeydown = function(e) {
                        if (e.keyCode === 13 && text.value !== "") {
                                //ws.send(now() + "\n" + "<" + name + "> " + text.value + "\n");
                                ws.send("<" + name + "> " + text.value + "\n");
                                text.value = "";
                                chat.scrollTop = 0;
                        }
                };
                $("#send").click(function() {
                        if (text.value !== "") {
                                ws.send("<" + name + "> " + text.value + "\n");
                                text.value = "";
                                chat.scrollTop = 0;
                        }
                });

/*              $("#refresh").click(function() {
                        var name = $("#username").text();
                        const form = document.createElement('form');
                        form.setAttribute('action', '/chat'); // 设置表单的action属性
                        form.setAttribute('method', 'post'); // 设置表单的method属性
                        form.setAttribute('enctype', 'multipart/form-data'); // 设置为带有文件的表单类型

                        // 创建一个input元素来输入用户名
                        const usernameInput = document.createElement('input');
                        usernameInput.setAttribute('type', 'text');
                        usernameInput.setAttribute('name', 'username');
                        usernameInput.setAttribute('value', name); // 设置要提交的用户名
                        form.appendChild(usernameInput); // 添加到表单中

                        // 将表单添加到页面并提交
                        document.body.appendChild(form);
                        form.submit();

                        // 可选：提交后移除表单
                        document.body.removeChild(form);
                        //window.location.href = "/chat_login";
                }); */
                $("#refresh").click(function() {
                        var uname = $("#username").text();
                        var rname = $("#roomname").text();
                        window.location.href = "/chat/"+rname+"?username="+uname;
                });

                $("#quit").click(function() {
                        window.location.href = "/chat_room";
                });

                console.log("init finished");
        }


        function createWebSocket() {
                try {
                        console.log("create socket begin");
                var roomname = $("#roomname").text();
                var url = socketPrefix + "://" + window.location.host + "/wschat/" + roomname;
                        if (!flagE){
                                ws = new WebSocket(url);
                                init();
                                flagE = true;
                        }
                        console.log("create socket end");
                } catch(e) {
                        console.log('catch' + e);
                        createWebSocket();
                        //reconnect();
                }
        }
/*
        document.addEventListener('touchmove', function (e) {
                if (e.touches.length > 1) {
                        return;
                }
                //e.preventDefault();
        }, {passive: false})
        */
/*
        function reconnect() {
                console.log("reconnect start");
                var tt;
                if(lockReconnect) {
                        return;
                };
                lockReconnect = true;
                //没连接上会一直重连，设置延迟避免请求过多，  有定时任务 先取消再设置
                tt && clearTimeout(tt);
                tt = setTimeout(function () {
                        createWebSocket();
                        lockReconnect = false;
                }, 3000);
                console.log("reconnect finished");
        }
*/
});

