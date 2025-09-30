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
      .catch((err) => {
        alert(err);
      });
  }
</script>

<main>
  <button class="btn" on:click={loadDotDfa}>Load .dfa</button>
</main>

<style>
  main {
    display: flexbox;
    align-items: center;
    justify-content: center;
  }

  .btn {
    width: 60%;
    height: 50%;
  }
</style>
