<script lang="ts">
  import { LoadDotDfa, LoadDotIn, EvaluateInput, SaveOutput } from "../wailsjs/go/main/App.js";

  let dfaContent: string = "";
  let dfaFilename: string = "";
  let inputLines: string[] = [];
  let inputFilename: string = "";
  let outputResults: boolean[] = [];
  let statusMessage: string = "";
  let previousDfaFilename: string = "";

  $: canProcess = dfaContent !== "" && inputLines.length > 0;

  let showFileTypeDialog = false;

  function loadFile(): void {
    showFileTypeDialog = true;
  }

  async function loadDfaFile(): Promise<void> {
    showFileTypeDialog = false;
    try {
      const result = await LoadDotDfa();
      if (!result) return; // User cancelled

      dfaContent = result.content;
      dfaFilename = result.filename;
      statusMessage = `DFA table from ${dfaFilename} has been successfully loaded.`;
      previousDfaFilename = dfaFilename;
    } catch (err) {
      if (previousDfaFilename) {
        statusMessage = `Unable to load DFA file due to invalid content. The program will be using the content from the most recently successfully loaded ${previousDfaFilename}.`;
      } else {
        statusMessage = `Unable to load DFA file due to invalid content.`;
      }
    }
  }

  async function loadInFile(): Promise<void> {
    showFileTypeDialog = false;
    try {
      const result = await LoadDotIn();
      if (!result || result.inputLines.length === 0) return; // User cancelled

      inputLines = result.inputLines;
      inputFilename = result.filename;
      statusMessage = `Input from file ${inputFilename} has been successfully loaded.`;
    } catch (err) {
      statusMessage = `Unable to load input file due to invalid content.`;
    }
  }

  async function processInput(): Promise<void> {
    if (!canProcess) return;

    let results: boolean[] = [];

    try {
      results = await EvaluateInput(inputLines);
      outputResults = results;

    } catch (err) {
      statusMessage = `Error processing input: ${err}`;
    }

    // Save the output to file
    try {
      await SaveOutput(results);
      statusMessage = `Input from ${inputFilename} successfully processed using DFA table from ${dfaFilename}. Output saved to ${inputFilename.replace('.in', '.out')}.`;
    } catch (saveErr) {
      statusMessage = `Input processed successfully, but failed to save output: ${saveErr}`;
    }
  }

  function parseTransitionTable(content: string): string[][] {
    if (!content) return [];
    const lines = content.split("\n");
    // Skip the first line (alphabet) and parse the rest
    return lines.slice(1).map(line => line.split(","));
  }
</script>

<main>
  <div class="controls">
    <button class="btn" on:click={loadFile}>Load File</button>
    <button class="btn" on:click={processInput} disabled={!canProcess}>Process</button>
  </div>

  {#if showFileTypeDialog}
    <div class="modal-overlay" on:click={() => showFileTypeDialog = false}>
      <div class="modal" on:click|stopPropagation>
        <h2>Select File Type</h2>
        <p>Which type of file would you like to load?</p>
        <div class="modal-buttons">
          <button class="btn modal-btn" on:click={loadDfaFile}>DFA File (.dfa)</button>
          <button class="btn modal-btn" on:click={loadInFile}>Input File (.in)</button>
          <button class="btn modal-btn cancel" on:click={() => showFileTypeDialog = false}>Cancel</button>
        </div>
      </div>
    </div>
  {/if}

  {#if statusMessage}
    <div class="status">{statusMessage}</div>
  {/if}

  <div class="displays">
    <div class="display-section">
      <h3>Transition Table</h3>
      <div class="display-content">
        {#if dfaContent}
          <table>
            <thead>
              <tr>
                <th></th>
                <th>State</th>
                <th>0</th>
                <th>1</th>
              </tr>
            </thead>
            <tbody>
              {#each parseTransitionTable(dfaContent) as row}
                <tr>
                  {#each row as cell}
                    <td style="color:white">{cell}</td>
                  {/each}
                </tr>
              {/each}
            </tbody>
          </table>
        {:else}
          <p class="empty">No DFA loaded</p>
        {/if}
      </div>
    </div>

    <div class="display-section">
      <h3>Input</h3>
      <div class="display-content">
        {#if inputLines.length > 0}
          <ul>
            {#each inputLines as line}
              <li>{line}</li>
            {/each}
          </ul>
        {:else}
          <p class="empty">No input loaded</p>
        {/if}
      </div>
    </div>

    <div class="display-section">
      <h3>Output</h3>
      <div class="display-content">
        {#if outputResults.length > 0}
          <ul>
            {#each outputResults as result}
              <li>{result ? "VALID" : "INVALID"}</li>
            {/each}
          </ul>
        {:else}
          <p class="empty">No output yet</p>
        {/if}
      </div>
    </div>
  </div>
</main>

<style>
  main {
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
    color: #fff;
  }

  .controls {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }

  .btn {
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 5px;
    transition: background-color 0.3s;
  }

  .btn:hover:not(:disabled) {
    background-color: #0056b3;
  }

  .btn:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }

  .status {
    padding: 12px;
    margin-bottom: 20px;
    background-color: #e7f3ff;
    border-left: 4px solid #007bff;
    border-radius: 4px;
    font-size: 14px;
    color: black;
  }

  .displays {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 20px;
  }

  .display-section {
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
  }

  .display-section h3 {
    margin: 0;
    padding: 12px;
    background-color: #f5f5f5;
    border-bottom: 1px solid #ddd;
    font-size: 16px;
    color: #000;
  }

  .display-content {
    padding: 12px;
    min-height: 150px;
    max-height: 400px;
    overflow-y: auto;
  }

  .empty {
    color: #999;
    font-style: italic;
    margin: 0;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th, td {
    padding: 8px;
    text-align: left;
    border-bottom: 1px solid #ddd;
  }

  th {
    background-color: #f5f5f5;
    font-weight: bold;
    color: black;
  }

  ul {
    list-style: none;
    padding: 0;
    margin: 0;
  }

  li {
    padding: 6px 0;
    border-bottom: 1px solid #eee;
  }

  li:last-child {
    border-bottom: none;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: white;
    padding: 30px;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    max-width: 400px;
    width: 90%;
  }

  .modal h2 {
    margin-top: 0;
    margin-bottom: 10px;
  }

  .modal p {
    margin-bottom: 20px;
    color: #666;
  }

  .modal-buttons {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .modal-btn {
    width: 100%;
  }

  .modal-btn.cancel {
    background-color: #6c757d;
  }

  .modal-btn.cancel:hover {
    background-color: #5a6268;
  }
</style>
