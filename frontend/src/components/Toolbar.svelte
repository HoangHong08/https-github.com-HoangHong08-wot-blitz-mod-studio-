<script>
  import { createEventDispatcher } from 'svelte'

  const dispatch = createEventDispatcher()

  function handleOpen() {
    dispatch('open')
  }

  function handleSave() {
    dispatch('save')
  }
</script>

<div class="toolbar">
  <div class="toolbar-left">
    <h1>WoT Blitz Mod Studio</h1>
  </div>
  <div class="toolbar-center">
    <span class="status">Ready</span>
  </div>
  <div class="toolbar-right">
    <button on:click={handleOpen} title="Open file">
      📁 Open
    </button>
    <button on:click={handleSave} title="Save file">
      💾 Save
    </button>
  </div>
</div>

<style>
  .toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.75rem 1rem;
    background: #2d2d2d;
    border-bottom: 1px solid #3e3e3e;
    height: 60px;
    gap: 1rem;
  }

  .toolbar-left {
    flex: 0 0 auto;
  }

  .toolbar-left h1 {
    margin: 0;
    font-size: 1.2rem;
    color: #4CAF50;
    font-weight: 600;
  }

  .toolbar-center {
    flex: 1;
    text-align: center;
  }

  .status {
    font-size: 0.9rem;
    color: #888;
  }

  .toolbar-right {
    display: flex;
    gap: 0.5rem;
    flex: 0 0 auto;
  }

  button {
    padding: 0.5rem 1rem;
    background: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background 0.2s;
  }

  button:hover {
    background: #45a049;
  }

  button:active {
    background: #3d8b40;
  }
</style>
