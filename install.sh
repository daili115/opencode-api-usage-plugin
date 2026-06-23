#!/bin/bash
# OpenCode API Usage Plugin Installer
# Usage: curl -sSL https://raw.githubusercontent.com/daili115/opencode-api-usage-plugin/main/install.sh | bash

set -e

REPO_URL="https://github.com/daili115/opencode-api-usage-plugin"
RAW_URL="https://raw.githubusercontent.com/daili115/opencode-api-usage-plugin/main"
PLUGIN_FILE="api_usage.go"
TARGET_DIR="internal/tui/components/chat"
EDITOR_FILE="internal/tui/components/chat/editor.go"

echo "🚀 OpenCode API Usage Plugin Installer"
echo "======================================"
echo ""

# Check if we're in an OpenCode repository
if [ ! -f "go.mod" ]; then
    echo "❌ Error: No go.mod found. Please run this script from the OpenCode repository root."
    exit 1
fi

if ! grep -q "github.com/opencode-ai/opencode" go.mod 2>/dev/null; then
    echo "⚠️  Warning: This doesn't appear to be an OpenCode repository."
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Download plugin file
echo "📥 Downloading plugin file..."
mkdir -p "$TARGET_DIR"
curl -sSL "$RAW_URL/api_usage.go" -o "$TARGET_DIR/$PLUGIN_FILE"

if [ ! -f "$TARGET_DIR/$PLUGIN_FILE" ]; then
    echo "❌ Error: Failed to download plugin file."
    exit 1
fi

echo "✅ Plugin file downloaded to $TARGET_DIR/$PLUGIN_FILE"
echo ""

# Check if editor.go exists
if [ ! -f "$EDITOR_FILE" ]; then
    echo "❌ Error: $EDITOR_FILE not found."
    exit 1
fi

# Backup editor.go
echo "💾 Backing up $EDITOR_FILE..."
cp "$EDITOR_FILE" "$EDITOR_FILE.bak"
echo "✅ Backup created: $EDITOR_FILE.bak"
echo ""

# Modify editor.go
echo "🔧 Modifying $EDITOR_FILE..."

# Check if already modified
if grep -q "apiUsage" "$EDITOR_FILE"; then
    echo "⚠️  Warning: Plugin appears to be already installed."
    read -p "Reinstall? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 0
    fi
fi

# Add apiUsage field to editorCmp struct
if ! grep -q "apiUsage.*\*apiUsageCmp" "$EDITOR_FILE"; then
    sed -i.bak 's/type editorCmp struct {/type editorCmp struct {\n\tapiUsage    *apiUsageCmp/' "$EDITOR_FILE"
fi

# Add apiUsage initialization in NewEditorCmp
if ! grep -q "apiUsage: NewAPIUsageCmp" "$EDITOR_FILE"; then
    sed -i.bak 's/return \&editorCmp{/return \&editorCmp{\n\t\tapiUsage: NewAPIUsageCmp().(*apiUsageCmp),/' "$EDITOR_FILE"
fi

# Add APIUsageMsg handling in Update
if ! grep -q "case APIUsageMsg:" "$EDITOR_FILE"; then
    sed -i.bak 's/case dialog.AttachmentAddedMsg:/case APIUsageMsg:\n\t\t_, cmd := m.apiUsage.Update(msg)\n\t\treturn m, cmd\n\tcase dialog.AttachmentAddedMsg:/' "$EDITOR_FILE"
fi

# Modify View method to include apiUsage
if ! grep -q "apiUsageView := m.apiUsage.View()" "$EDITOR_FILE"; then
    # This is a simplified modification - in reality you'd need more sophisticated parsing
    echo "⚠️  Please manually modify the View() method to include apiUsage rendering."
    echo "   Add the following before the final return statement:"
    echo ""
    echo "   apiUsageView := m.apiUsage.View()"
    echo "   if apiUsageView != \"\" {""
    echo "       return lipgloss.JoinVertical(lipgloss.Top, editorView, apiUsageView)"
    echo "   }"
    echo ""
fi

# Modify SetSize to reserve space for apiUsage
if ! grep -q "m.apiUsage.SetSize" "$EDITOR_FILE"; then
    sed -i.bak 's/m.textarea.SetHeight(height)/m.textarea.SetHeight(height - 3)\n\tm.apiUsage.SetSize(width, 3)/' "$EDITOR_FILE"
fi

# Add apiUsage bindings
if ! grep -q "m.apiUsage.BindingKeys()" "$EDITOR_FILE"; then
    sed -i.bak 's/bindings = append(bindings, layout.KeyMapToSlice(DeleteKeyMaps)...)/bindings = append(bindings, layout.KeyMapToSlice(DeleteKeyMaps)...)\n\tbindings = append(bindings, m.apiUsage.BindingKeys()...)/' "$EDITOR_FILE"
fi

echo "✅ Modifications applied to $EDITOR_FILE"
echo ""

# Verify compilation
echo "🔨 Verifying compilation..."
if go build ./...; then
    echo "✅ Compilation successful!"
    echo ""
    echo "🎉 Installation complete!"
    echo ""
    echo "Usage:"
    echo "  - The API usage stats will appear below the input box"
    echo "  - Press Ctrl+U to toggle display"
    echo ""
    echo "To uninstall:"
    echo "  cp $EDITOR_FILE.bak $EDITOR_FILE"
    echo "  rm $TARGET_DIR/$PLUGIN_FILE"
else
    echo "❌ Compilation failed. Please check the error messages above."
    echo ""
    echo "To restore the backup:"
    echo "  cp $EDITOR_FILE.bak $EDITOR_FILE"
    exit 1
fi
