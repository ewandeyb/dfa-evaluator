<script lang="ts">
  import { LoadDotDfa, LoadDotIn, EvaluateInput } from "../wailsjs/go/main/App.js";

  function loadDotDfa(): void {
    LoadDotDfa()
      .then((filename) => {
        alert(filename);
      })
      .catch((err) => {
        alert(err);
      });
  }

  let inputLines: string[] = [];

  function loadDotIn(): void {
    LoadDotIn()
      .then((ret) => {
        alert(ret);
        inputLines = ret;
      })
      .catch((err) => {
        alert(err);
      });
  }

  function evaluateInput(): void {
    EvaluateInput(inputLines)
      .then((ret) => {
        alert(ret);
      })
      .catch((err) => {
        alert(err);
      });
  }
</script>

<main>
  <button class="btn" on:click={loadDotDfa}>Load .dfa</button>
  <button class="btn" on:click={loadDotIn}>Load .in</button>
  <button class="btn" on:click={evaluateInput}>Evaluate Input</button>
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
