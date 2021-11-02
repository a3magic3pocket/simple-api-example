// Convert es6 to es5 : https://es6console.com/
let preventSubmit = true;

function validate(userName, password) {
  const errors = { username: "", password: "" };

  if (userName.length < 8) {
    errors["username"] = "8자 이상의 ID를 사용해야합니다.";
  } else if (!/[A-z0-9-_]+/.test(userName)) {
    errors["username"] = "ID는 (알파벳, 숫자, -, _)만 사용할 수 있습니다.";
  }

  if (password.length < 8) {
    errors["password"] = "8자 이상의 Password를 사용해야합니다.";
  }

  return errors;
}

function renderErrors(errors) {
  for (const key in errors) {
    const elem = document.querySelector(`#login-${key}-error`);
    if (typeof elem !== "undefined") {
      elem.innerText = errors[key];
    }
  }
}

function handleSubmit(event) {
  if (preventSubmit) {
    event.preventDefault();
  }

  const formElem = document.querySelector("form");
  const formData = new FormData(formElem);
  const userName = formData.get("UserName");
  const password = formData.get("Password");

  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  const redirectURL = urlParams.get("redirect-url");

  const errors = validate(userName, password);
  renderErrors(errors);
  for (const errorKey in errors) {
    if (errors[errorKey] !== "") {
      return;
    }
  }

  const fetchLoginAsync = async (url) => {
    let response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        UserName: formData.get("UserName"),
        Password: formData.get("Password"),
      }),
    });
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.blob();
  };

  const loginErrorElem = document.querySelector("#login-error");
  fetchLoginAsync("login")
    .then(() => {
      if (typeof loginErrorElem !== "undefined") {
        loginErrorElem.innerText = "";
      }
      formElem.action = "//" + redirectURL;
      console.log("formElem.action", formElem.action);
      preventSubmit = false;
      formElem.submit();
    })
    .catch((e) => {
      if (typeof loginErrorElem !== "undefined") {
        loginErrorElem.innerText = "ID 또는 Password를 잘못 입력하셨습니다.";
      }
    });
}
