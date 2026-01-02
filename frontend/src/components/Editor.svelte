<script>
  import { createEventDispatcher, onMount } from 'svelte'

  export let content = ''

  const dispatch = createEventDispatcher()

  let editorDiv
  let editor = null

  onMount(async () => {
    // In production, would use Monaco editor
    // For now, create a simple textarea for the prototype
    if (editorDiv) {
      editorDiv.innerHTML = ''
      const textarea = document.createElement('textarea')
      textarea.value = content
      textarea.style.width = '100%'
      textarea.style.height = '100%'
      textarea.style.padding = '1rem'
      textarea.style.fontFamily = 'monospace'
      textarea.style.fontSize = '0.9rem'
      textarea.style.border = 'none'
      textarea.style.background = '#1e1e1e'
      textarea.style.color = '#d4d4d4'
      textarea.style.resize = 'none'
      textarea.style.overflow = 'auto'

      textarea.addEventListener('input', (e) => {
        dispatch('contentChanged', e.target.value)
      })

      editorDiv.appendChild(textarea)
    }
  })

  $: if (editorDiv && editorDiv.firstChild && content !== editorDiv.firstChild.value) {
    editorDiv.firstChild.value = content
  }
</script>

<div class="editor">
  <div class="editor-header">
    <span>YAML Editor</span>
  </div>
  <div class="editor-content" bind:this={editorDiv} />
</div>

<style>
  .editor {
    display: flex;
    flex-direction: column;
    background: #1e1e1e;
    border-right: 1px solid #3e3e3e;
    overflow: hidden;
  }

  .editor-header {
    padding: 0.5rem 1rem;
    background: #2d2d2d;
    border-bottom: 1px solid #3e3e3e;
    font-size: 0.9rem;
    color: #cccccc;
    font-weight: 500;
  }

  .editor-content {
    flex: 1;
    overflow: hidden;
  }

  .editor-content :global(textarea) {
    outline: none;
  }

  .editor-content :global(textarea:focus) {
    background: #252526 !important;
  }
</style>
