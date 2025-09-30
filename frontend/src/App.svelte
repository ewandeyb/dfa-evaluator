<script lang="ts">
  import logo from "./assets/images/logo-universal.png";
  import { LoadDotDfa } from "../wailsjs/go/main/App.js";

  let resultText: string = "Please enter your name below 👇";
  let name: string;

  function loadDotDfa(): void {
    LoadDotDfa()
      .then((filename) => {
        alert(filename);
      })
      .catch(() => {});
  }
</script>

<main>
  <img alt="Wails logo" id="logo" src={logo} />
  <div class="result" id="result">{resultText}</div>
  <div class="input-box" id="input">
    <button class="btn" on:click={loadDotDfa}>Load .dfa</button>
  </div>
</main>
