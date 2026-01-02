<script>
  import Toolbar from './components/Toolbar.svelte'
  import Sidebar from './components/Sidebar.svelte'
  import Editor from './components/Editor.svelte'
  import Preview from './components/Preview.svelte'
  import { writable } from 'svelte/store'
  import { OpenFile, SaveFile, ParseYAML } from '../wailsjs/go/main/App.js'

  let currentFile = null
  let currentContent = writable('')
  let currentPackage = writable(null)
  let currentAssets = writable([])
  let fileIsDVPL = false
  let selectedControl = writable(null)

  async function handleOpenFile() {
    try {
      // In a real app, we'd use a file picker
      // For now, this is a placeholder
      console.log('Open file handler needed')
    } catch (error) {
      console.error('Error opening file:', error)
    }
  }

  async function handleSaveFile() {
    try {
      if (!currentFile) {
        alert('No file open')
        return
      }

      const content = $currentContent
      await SaveFile(currentFile, content, fileIsDVPL)
      alert('File saved successfully')
    } catch (error) {
      console.error('Error saving file:', error)
      alert('Error saving file: ' + error)
    }
  }

  function handleContentChanged(event) {
    currentContent.set(event.detail)
  }

  function handleControlSelected(event) {
    selectedControl.set(event.detail)
  }
</script>

<div class="app">
  <Toolbar 
    on:open={handleOpenFile}
    on:save={handleSaveFile}
  />
  
  <div class="main">
    <Sidebar 
      uiPackage={$currentPackage}
      on:selectControl={handleControlSelected}
    />
    <Editor 
      content={$currentContent}
      on:contentChanged={handleContentChanged}
    />
    <Preview 
      uiPackage={$currentPackage}
      selectedControl={$selectedControl}
    />
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    background: #1e1e1e;
    color: #e0e0e0;
  }

  :global(*) {
    box-sizing: border-box;
  }

  .app {
    display: flex;
    flex-direction: column;
    height: 100vh;
    width: 100vw;
    overflow: hidden;
  }

  .main {
    display: grid;
    grid-template-columns: 300px 1fr 500px;
    flex: 1;
    overflow: hidden;
    gap: 0;
  }

  @media (max-width: 1200px) {
    .main {
      grid-template-columns: 250px 1fr 400px;
    }
  }

  @media (max-width: 768px) {
    .main {
      grid-template-columns: 1fr;
    }
  }
</style>
