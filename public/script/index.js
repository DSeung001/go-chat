$(function () {
    const webSocketScheme = window.location.protocol == "https:" ? 'wss://' : 'ws://';
    const baseURI = window.location.hostname + (location.port ? ':' + location.port : '');

    // New WebSocket을 실행한 순간 go의 socketHandler가 실행됨
    const websocket = new WebSocket(webSocketScheme + baseURI + '/ws');
    let nameInput = $("#name-input");
    let userName

    // data에 따른 메시지 태그 추가 함수
    function log(data) {
        let message = ""

        switch (Number(data.Type)) {
            case 0 :
                message = `<li class="message notice-message">${data.Message}</li>`;
                break
            case 1 :
                if (data.Author === userName) {
                    message = `<li class="message user-message">${data.Author} : ${data.Message}</li>`
                } else {
                    message = `<li class="message other-message">${data.Author} : ${data.Message}</li>`
                }
                break
            default :
                message = `<li class="message notice-message">내부 문제가 발생했습니다.</li>`;
        }
        $('#messages').append(message)
    }

    // GET 요청으로 현재 참여중인 유저의 리스트를 가져오고 적용
    function getUsers() {
        $.get("/getUsers")
            .done(function (response) {
                $("#participants").empty()
                response.forEach(name =>{
                    if (name !== ""){
                        $("#participants").append(`<li>${name}</li>`)
                    }
                })
            })
            .fail(function (jqXHR, textStatus, errorThrown) {
                console.log(jqXHR)
                console.log("API 요청 실패:", textStatus, errorThrown);
            });
    }

    // websocket에 메시지가 왔을 때 실행
    websocket.onmessage = function (e) {
        getUsers()
        log(JSON.parse(e.data))
    };

    // websocket에 에러가 발생했을 때 실행
    websocket.onerror = function (e) {
        log('에러 발생');
        console.log(e);
    };

    // 이름 입력 시 실행
    $("#name-form").submit(function (e) {
        e.preventDefault()
        $.post("/login", {name: nameInput.val() })
            .done(function () {
                userName = nameInput.val()
                $("#name-save").hide()
                nameInput.prop("disabled", true);
                $(".chat-container").removeClass("hide");
            })
            .fail(function (jqXHR, textStatus, errorThrown) {
                console.log(jqXHR)
                console.log("API 요청 실패:", textStatus, errorThrown);
                alert("중복되는 닉네임입니다.")
            });
    })

    // 메시지 입력시 실행
    $('#chat-form').submit(function (e) {
        e.preventDefault();
        let data = $('#chat-text').val();
        if (data) {
            // websocket.send으로 연결중인 웹 소캣에 메세지 전송
            // 현재 연결중인 websocket connection 실행
            // socketHandler에서 계속 돌고 있는 go p.Read() 안에 peer.Conn.ReadMessage()에 메세지가 들어오고 다음 코드 실행
            websocket.send(
                JSON.stringify({
                    Author: userName,
                    Message: data
                }));
            window.scrollTo(0, document.body.scrollHeight)
            $('#chat-text').val('');
        }
    });
});