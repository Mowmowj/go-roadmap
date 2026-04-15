#!/usr/bin/env python3
"""
Go Roadmap 静态网站生成器
从 Go 源代码文件生成带有语法高亮的静态 HTML 网站
"""

import os
import re
import html
import json
import shutil

# 章节元数据
CHAPTERS = [
    {"dir": "01-getting-started",       "title": "Go 语言入门",         "emoji": "🚀", "level": "初级", "tag": "必须"},
    {"dir": "02-variables-constants",   "title": "变量与常量",          "emoji": "📦", "level": "初级", "tag": "必须"},
    {"dir": "03-data-types",            "title": "数据类型",            "emoji": "🔢", "level": "初级", "tag": "必须"},
    {"dir": "04-composite-types",       "title": "复合类型",            "emoji": "🧩", "level": "初级", "tag": "必须"},
    {"dir": "05-control-flow",          "title": "控制流",              "emoji": "🔀", "level": "初级", "tag": "必须"},
    {"dir": "06-functions",             "title": "函数",                "emoji": "⚡", "level": "初级", "tag": "必须"},
    {"dir": "07-pointers",              "title": "指针",                "emoji": "👉", "level": "中级", "tag": "必须"},
    {"dir": "08-methods-interfaces",    "title": "方法与接口",          "emoji": "🔌", "level": "中级", "tag": "必须"},
    {"dir": "09-generics",              "title": "泛型",                "emoji": "🧬", "level": "中级", "tag": "必须"},
    {"dir": "10-error-handling",        "title": "错误处理",            "emoji": "🛡️", "level": "中级", "tag": "必须"},
    {"dir": "11-modules-packages",      "title": "模块与包",            "emoji": "📚", "level": "中级", "tag": "必须"},
    {"dir": "12-concurrency",           "title": "并发编程",            "emoji": "🔄", "level": "高级", "tag": "必须"},
    {"dir": "13-standard-library",      "title": "标准库",              "emoji": "📖", "level": "中级", "tag": "必须"},
    {"dir": "14-testing",               "title": "测试",                "emoji": "🧪", "level": "中级", "tag": "必须"},
    {"dir": "15-ecosystem",             "title": "生态系统",            "emoji": "🌍", "level": "中级", "tag": "推荐"},
    {"dir": "16-toolchain",             "title": "工具链",              "emoji": "🔧", "level": "中级", "tag": "推荐"},
    {"dir": "17-advanced",              "title": "高级主题",            "emoji": "🎯", "level": "高级", "tag": "可选"},
]

PHASES = [
    {"name": "语言基础",           "range": (1, 6),  "color": "#10b981", "level": "🟢 初级"},
    {"name": "核心进阶",           "range": (7, 11), "color": "#f59e0b", "level": "🟡 中级"},
    {"name": "并发编程",           "range": (12, 12),"color": "#ef4444", "level": "🔴 高级"},
    {"name": "标准库与工程实践",   "range": (13, 14),"color": "#f59e0b", "level": "🟡 中级"},
    {"name": "生态与工具",         "range": (15, 16),"color": "#f59e0b", "level": "🟡 中级"},
    {"name": "高级主题",           "range": (17, 17),"color": "#ef4444", "level": "🔴 高级"},
]

def get_level_color(level):
    colors = {"初级": "#10b981", "中级": "#f59e0b", "高级": "#ef4444"}
    return colors.get(level, "#6b7280")

def get_tag_color(tag):
    colors = {"必须": "#ef4444", "推荐": "#3b82f6", "可选": "#8b5cf6"}
    return colors.get(tag, "#6b7280")

def extract_sections(code):
    """从 Go 代码中提取章节标题"""
    sections = []
    for line in code.split('\n'):
        # 匹配形如 // 1. xxx 或 // 📚 概念：xxx 的行
        m = re.match(r'^\s*//\s+(\d+)\.\s+(.+)', line)
        if m:
            num, title = m.group(1), m.group(2).strip()
            if not any(c in title for c in ['=', '-', '+']):
                sections.append({"id": f"section-{num}", "title": f"{num}. {title}"})
    return sections

def read_go_files(root, chapter_dir):
    """读取章节目录下的所有 .go 文件"""
    files = []
    chapter_path = os.path.join(root, chapter_dir)
    if not os.path.isdir(chapter_path):
        return files
    for fname in sorted(os.listdir(chapter_path)):
        if fname.endswith('.go'):
            fpath = os.path.join(chapter_path, fname)
            with open(fpath, 'r', encoding='utf-8') as f:
                files.append({"name": fname, "content": f.read()})
    return files

def generate_html_page(chapter, files, chapters, root):
    """生成单个章节的 HTML 页面"""
    idx = chapters.index(chapter)
    prev_ch = chapters[idx - 1] if idx > 0 else None
    next_ch = chapters[idx + 1] if idx < len(chapters) - 1 else None
    
    code_blocks = ""
    for f in files:
        escaped = html.escape(f["content"])
        code_blocks += f"""
        <div class="file-block">
            <div class="file-header">
                <span class="file-icon">📄</span>
                <span class="file-name">{html.escape(f['name'])}</span>
                <button class="copy-btn" onclick="copyCode(this)" title="复制代码">
                    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"/><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"/></svg>
                </button>
            </div>
            <pre><code class="language-go">{escaped}</code></pre>
        </div>
"""

    nav_links = ""
    for phase in PHASES:
        nav_links += f'<div class="nav-phase" style="border-left-color: {phase["color"]}">'
        nav_links += f'<div class="nav-phase-title">{phase["level"]} {phase["name"]}</div>'
        for ch in chapters:
            ch_num = int(ch["dir"].split("-")[0])
            if phase["range"][0] <= ch_num <= phase["range"][1]:
                active = "active" if ch == chapter else ""
                nav_links += f'<a href="{ch["dir"]}.html" class="nav-item {active}">'
                nav_links += f'<span class="nav-emoji">{ch["emoji"]}</span>'
                nav_links += f'<span class="nav-text">{ch["title"]}</span>'
                nav_links += f'</a>'
        nav_links += '</div>'

    prev_link = f'<a href="{prev_ch["dir"]}.html" class="nav-btn prev-btn">← {prev_ch["emoji"]} {prev_ch["title"]}</a>' if prev_ch else '<span></span>'
    next_link = f'<a href="{next_ch["dir"]}.html" class="nav-btn next-btn">{next_ch["emoji"]} {next_ch["title"]} →</a>' if next_ch else '<span></span>'

    level_color = get_level_color(chapter["level"])
    tag_color = get_tag_color(chapter["tag"])

    return f"""<!DOCTYPE html>
<html lang="zh-CN" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{chapter["emoji"]} {chapter["title"]} - Go 学习路线图</title>
    <link rel="stylesheet" href="assets/style.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/languages/go.min.js"></script>
</head>
<body>
    <button class="sidebar-toggle" onclick="toggleSidebar()" aria-label="切换侧边栏">☰</button>
    <nav class="sidebar" id="sidebar">
        <div class="sidebar-header">
            <a href="index.html" class="logo">
                <span class="logo-icon">🚀</span>
                <span class="logo-text">Go 路线图</span>
            </a>
        </div>
        <div class="sidebar-nav">
            {nav_links}
        </div>
    </nav>
    <main class="content">
        <div class="chapter-header">
            <div class="chapter-meta">
                <span class="badge" style="background: {level_color}">{chapter["level"]}</span>
                <span class="badge" style="background: {tag_color}">{chapter["tag"]}</span>
            </div>
            <h1>{chapter["emoji"]} {chapter["title"]}</h1>
        </div>
        <div class="code-content">
            {code_blocks}
        </div>
        <div class="page-nav">
            {prev_link}
            {next_link}
        </div>
    </main>
    <script>
        hljs.highlightAll();
        function copyCode(btn) {{
            const code = btn.closest('.file-block').querySelector('code').textContent;
            navigator.clipboard.writeText(code).then(() => {{
                btn.innerHTML = '✓';
                setTimeout(() => {{
                    btn.innerHTML = '<svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 010 1.5h-1.5a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-1.5a.75.75 0 011.5 0v1.5A1.75 1.75 0 019.25 16h-7.5A1.75 1.75 0 010 14.25v-7.5z"/><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0114.25 11h-7.5A1.75 1.75 0 015 9.25v-7.5zm1.75-.25a.25.25 0 00-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 00.25-.25v-7.5a.25.25 0 00-.25-.25h-7.5z"/></svg>';
                }}, 1500);
            }});
        }}
        function toggleSidebar() {{
            document.getElementById('sidebar').classList.toggle('open');
        }}
        document.querySelector('.content').addEventListener('click', () => {{
            document.getElementById('sidebar').classList.remove('open');
        }});
    </script>
</body>
</html>"""


def generate_index(chapters):
    """生成首页"""
    phase_cards = ""
    for phase in PHASES:
        items = ""
        for ch in chapters:
            ch_num = int(ch["dir"].split("-")[0])
            if phase["range"][0] <= ch_num <= phase["range"][1]:
                level_color = get_level_color(ch["level"])
                tag_color = get_tag_color(ch["tag"])
                items += f"""
                <a href="{ch['dir']}.html" class="chapter-card">
                    <span class="chapter-emoji">{ch['emoji']}</span>
                    <div class="chapter-info">
                        <span class="chapter-title">{ch['dir'].split('-')[0]}. {ch['title']}</span>
                        <div class="chapter-badges">
                            <span class="badge-sm" style="background: {level_color}">{ch['level']}</span>
                            <span class="badge-sm" style="background: {tag_color}">{ch['tag']}</span>
                        </div>
                    </div>
                    <span class="chapter-arrow">→</span>
                </a>"""
        
        phase_cards += f"""
        <div class="phase-section">
            <div class="phase-header" style="border-left-color: {phase['color']}">
                <span class="phase-level">{phase['level']}</span>
                <h2>{phase['name']}</h2>
            </div>
            <div class="chapter-list">
                {items}
            </div>
        </div>"""

    nav_links = ""
    for phase in PHASES:
        nav_links += f'<div class="nav-phase" style="border-left-color: {phase["color"]}">'
        nav_links += f'<div class="nav-phase-title">{phase["level"]} {phase["name"]}</div>'
        for ch in chapters:
            ch_num = int(ch["dir"].split("-")[0])
            if phase["range"][0] <= ch_num <= phase["range"][1]:
                nav_links += f'<a href="{ch["dir"]}.html" class="nav-item">'
                nav_links += f'<span class="nav-emoji">{ch["emoji"]}</span>'
                nav_links += f'<span class="nav-text">{ch["title"]}</span>'
                nav_links += f'</a>'
        nav_links += '</div>'

    return f"""<!DOCTYPE html>
<html lang="zh-CN" data-theme="dark">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go 语言完整学习路线图</title>
    <link rel="stylesheet" href="assets/style.css">
</head>
<body>
    <button class="sidebar-toggle" onclick="toggleSidebar()" aria-label="切换侧边栏">☰</button>
    <nav class="sidebar" id="sidebar">
        <div class="sidebar-header">
            <a href="index.html" class="logo">
                <span class="logo-icon">🚀</span>
                <span class="logo-text">Go 路线图</span>
            </a>
        </div>
        <div class="sidebar-nav">
            {nav_links}
        </div>
    </nav>
    <main class="content">
        <div class="hero">
            <h1>🚀 Go 语言完整学习路线图</h1>
            <p class="hero-sub">从零到精通 Go 语言的系统化学习项目，所有代码可直接运行</p>
            <div class="hero-badges">
                <span class="hero-badge green">🟢 初级</span>
                <span class="hero-badge yellow">🟡 中级</span>
                <span class="hero-badge red">🔴 高级</span>
            </div>
            <p class="hero-ref">基于 <a href="https://roadmap.sh/golang" target="_blank">roadmap.sh/golang</a></p>
        </div>

        <div class="legend">
            <div class="legend-item"><span class="legend-dot" style="background:#10b981"></span> 初级 - 零基础必学</div>
            <div class="legend-item"><span class="legend-dot" style="background:#f59e0b"></span> 中级 - 日常开发必备</div>
            <div class="legend-item"><span class="legend-dot" style="background:#ef4444"></span> 高级 - 深入理解优化</div>
            <div class="legend-sep"></div>
            <div class="legend-item"><span class="legend-dot" style="background:#ef4444"></span> 必须 - 不学无法继续</div>
            <div class="legend-item"><span class="legend-dot" style="background:#3b82f6"></span> 推荐 - 强烈建议掌握</div>
            <div class="legend-item"><span class="legend-dot" style="background:#8b5cf6"></span> 可选 - 按方向选学</div>
        </div>

        {phase_cards}

        <div class="how-to-run">
            <h2>🏃 如何运行</h2>
            <pre><code># 克隆项目
git clone &lt;repo-url&gt;
cd go-roadmap

# 运行某个章节
cd 01-getting-started && go run main.go

# 运行测试
cd 14-testing && go test -v ./...</code></pre>
        </div>

        <div class="tips">
            <h2>📖 学习建议</h2>
            <div class="tips-grid">
                <div class="tip-card">
                    <span class="tip-icon">📋</span>
                    <strong>按顺序学习</strong>
                    <p>章节间有依赖关系，建议从 01 开始</p>
                </div>
                <div class="tip-card">
                    <span class="tip-icon">💻</span>
                    <strong>动手实践</strong>
                    <p>每个文件的代码都可运行，修改后观察变化</p>
                </div>
                <div class="tip-card">
                    <span class="tip-icon">⭐</span>
                    <strong>先必须后可选</strong>
                    <p>优先完成必须标记内容，再学可选内容</p>
                </div>
                <div class="tip-card">
                    <span class="tip-icon">📝</span>
                    <strong>做练习</strong>
                    <p>每章末尾都有练习，动手完成它们</p>
                </div>
            </div>
        </div>
    </main>
    <script>
        function toggleSidebar() {{
            document.getElementById('sidebar').classList.toggle('open');
        }}
        document.querySelector('.content').addEventListener('click', () => {{
            document.getElementById('sidebar').classList.remove('open');
        }});
    </script>
</body>
</html>"""


def generate_css():
    return """/* ============================================
   Go Roadmap - Dark Theme Site Styles
   ============================================ */

:root {
    --bg-primary: #0d1117;
    --bg-secondary: #161b22;
    --bg-tertiary: #21262d;
    --bg-hover: #30363d;
    --border: #30363d;
    --text-primary: #e6edf3;
    --text-secondary: #8b949e;
    --text-muted: #6e7681;
    --accent: #58a6ff;
    --accent-subtle: #1f6feb;
    --green: #10b981;
    --yellow: #f59e0b;
    --red: #ef4444;
    --sidebar-width: 280px;
    --header-height: 60px;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

html {
    scroll-behavior: smooth;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Noto Sans', Helvetica, Arial, sans-serif;
    background: var(--bg-primary);
    color: var(--text-primary);
    line-height: 1.6;
    display: flex;
    min-height: 100vh;
}

/* Sidebar */
.sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: var(--sidebar-width);
    background: var(--bg-secondary);
    border-right: 1px solid var(--border);
    overflow-y: auto;
    z-index: 100;
    transition: transform 0.3s ease;
}

.sidebar-header {
    padding: 16px 20px;
    border-bottom: 1px solid var(--border);
    position: sticky;
    top: 0;
    background: var(--bg-secondary);
    z-index: 1;
}

.logo {
    display: flex;
    align-items: center;
    gap: 8px;
    text-decoration: none;
    color: var(--text-primary);
    font-weight: 600;
    font-size: 18px;
}

.logo-icon { font-size: 24px; }

.sidebar-nav {
    padding: 12px 0;
}

.nav-phase {
    border-left: 3px solid transparent;
    margin: 4px 0;
    padding-left: 0;
}

.nav-phase-title {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--text-muted);
    padding: 8px 20px 4px;
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 20px;
    text-decoration: none;
    color: var(--text-secondary);
    font-size: 14px;
    transition: all 0.15s ease;
    border-left: 3px solid transparent;
    margin-left: -3px;
}

.nav-item:hover {
    color: var(--text-primary);
    background: var(--bg-tertiary);
}

.nav-item.active {
    color: var(--accent);
    background: rgba(88, 166, 255, 0.1);
    border-left-color: var(--accent);
}

.nav-emoji { font-size: 16px; flex-shrink: 0; }
.nav-text { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* Sidebar toggle button */
.sidebar-toggle {
    display: none;
    position: fixed;
    top: 12px;
    left: 12px;
    z-index: 200;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    color: var(--text-primary);
    font-size: 20px;
    padding: 6px 10px;
    border-radius: 8px;
    cursor: pointer;
}

/* Main Content */
.content {
    flex: 1;
    margin-left: var(--sidebar-width);
    padding: 40px;
    max-width: 960px;
    min-height: 100vh;
}

/* Chapter Header */
.chapter-header {
    margin-bottom: 32px;
}

.chapter-meta {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
}

.badge {
    display: inline-block;
    padding: 2px 10px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 600;
    color: white;
}

.chapter-header h1 {
    font-size: 32px;
    font-weight: 700;
    line-height: 1.3;
}

/* Code blocks */
.file-block {
    margin-bottom: 24px;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid var(--border);
}

.file-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    background: var(--bg-tertiary);
    border-bottom: 1px solid var(--border);
    font-size: 13px;
    color: var(--text-secondary);
}

.file-name { font-family: 'SF Mono', 'Fira Code', monospace; }
.file-icon { font-size: 14px; }

.copy-btn {
    margin-left: auto;
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    transition: color 0.15s;
}

.copy-btn:hover { color: var(--text-primary); }

.file-block pre {
    margin: 0;
    padding: 0;
}

.file-block pre code {
    display: block;
    padding: 16px !important;
    font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
    font-size: 13px;
    line-height: 1.6;
    overflow-x: auto;
    background: var(--bg-primary) !important;
    tab-size: 4;
}

/* Page navigation */
.page-nav {
    display: flex;
    justify-content: space-between;
    margin-top: 48px;
    padding-top: 24px;
    border-top: 1px solid var(--border);
}

.nav-btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 8px;
    text-decoration: none;
    color: var(--text-primary);
    font-size: 14px;
    transition: all 0.15s;
}

.nav-btn:hover {
    background: var(--bg-tertiary);
    border-color: var(--accent);
}

/* ============================================
   Index / Home Page
   ============================================ */

.hero {
    text-align: center;
    padding: 48px 0 32px;
}

.hero h1 {
    font-size: 36px;
    font-weight: 800;
    margin-bottom: 12px;
    background: linear-gradient(135deg, #58a6ff, #10b981);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.hero-sub {
    color: var(--text-secondary);
    font-size: 18px;
    margin-bottom: 16px;
}

.hero-badges {
    display: flex;
    justify-content: center;
    gap: 12px;
    margin-bottom: 12px;
}

.hero-badge {
    padding: 4px 14px;
    border-radius: 16px;
    font-size: 13px;
    font-weight: 600;
}

.hero-badge.green { background: rgba(16,185,129,0.15); color: #10b981; }
.hero-badge.yellow { background: rgba(245,158,11,0.15); color: #f59e0b; }
.hero-badge.red { background: rgba(239,68,68,0.15); color: #ef4444; }

.hero-ref {
    color: var(--text-muted);
    font-size: 14px;
}

.hero-ref a {
    color: var(--accent);
    text-decoration: none;
}

.hero-ref a:hover { text-decoration: underline; }

/* Legend */
.legend {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    padding: 16px 20px;
    background: var(--bg-secondary);
    border-radius: 8px;
    border: 1px solid var(--border);
    margin-bottom: 32px;
    align-items: center;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: var(--text-secondary);
}

.legend-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
}

.legend-sep {
    width: 1px;
    height: 20px;
    background: var(--border);
}

/* Phase sections */
.phase-section {
    margin-bottom: 32px;
}

.phase-header {
    border-left: 4px solid;
    padding: 8px 16px;
    margin-bottom: 12px;
}

.phase-level {
    font-size: 12px;
    color: var(--text-muted);
}

.phase-header h2 {
    font-size: 20px;
    font-weight: 700;
}

.chapter-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.chapter-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 8px;
    text-decoration: none;
    color: var(--text-primary);
    transition: all 0.15s;
}

.chapter-card:hover {
    background: var(--bg-tertiary);
    border-color: var(--accent);
    transform: translateX(4px);
}

.chapter-emoji { font-size: 24px; flex-shrink: 0; }

.chapter-info { flex: 1; }

.chapter-title {
    font-weight: 600;
    display: block;
    margin-bottom: 2px;
}

.chapter-badges {
    display: flex;
    gap: 6px;
}

.badge-sm {
    padding: 1px 8px;
    border-radius: 8px;
    font-size: 11px;
    font-weight: 600;
    color: white;
}

.chapter-arrow {
    color: var(--text-muted);
    font-size: 18px;
    transition: transform 0.15s;
}

.chapter-card:hover .chapter-arrow {
    transform: translateX(4px);
    color: var(--accent);
}

/* How to run */
.how-to-run {
    margin-top: 32px;
    padding: 24px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 8px;
}

.how-to-run h2 {
    margin-bottom: 16px;
    font-size: 20px;
}

.how-to-run pre {
    background: var(--bg-primary);
    border: 1px solid var(--border);
    border-radius: 6px;
    padding: 16px;
    overflow-x: auto;
}

.how-to-run code {
    font-family: 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
    line-height: 1.6;
    color: var(--text-primary);
}

/* Tips */
.tips {
    margin-top: 24px;
}

.tips h2 {
    margin-bottom: 16px;
    font-size: 20px;
}

.tips-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 12px;
}

.tip-card {
    padding: 16px;
    background: var(--bg-secondary);
    border: 1px solid var(--border);
    border-radius: 8px;
    text-align: center;
}

.tip-icon { font-size: 28px; display: block; margin-bottom: 8px; }

.tip-card strong {
    display: block;
    margin-bottom: 4px;
    color: var(--text-primary);
}

.tip-card p {
    font-size: 13px;
    color: var(--text-secondary);
}

/* ============================================
   Responsive
   ============================================ */

@media (max-width: 768px) {
    .sidebar {
        transform: translateX(-100%);
    }

    .sidebar.open {
        transform: translateX(0);
        box-shadow: 4px 0 20px rgba(0,0,0,0.5);
    }

    .sidebar-toggle {
        display: block;
    }

    .content {
        margin-left: 0;
        padding: 20px 16px;
        padding-top: 60px;
    }

    .hero h1 { font-size: 24px; }
    .hero-sub { font-size: 15px; }
    .chapter-header h1 { font-size: 24px; }
    .page-nav { flex-direction: column; gap: 8px; }
    .nav-btn { justify-content: center; }
    .legend { flex-direction: column; gap: 8px; }
    .legend-sep { width: 100%; height: 1px; }
}

@media (max-width: 480px) {
    .tips-grid { grid-template-columns: 1fr; }
    .hero-badges { flex-direction: column; align-items: center; }
}

/* Scrollbar styling */
::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

::-webkit-scrollbar-track {
    background: var(--bg-primary);
}

::-webkit-scrollbar-thumb {
    background: var(--bg-tertiary);
    border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
    background: var(--bg-hover);
}
"""

def build_site(root_dir, output_dir):
    """构建整个站点"""
    # 清理输出目录
    if os.path.exists(output_dir):
        shutil.rmtree(output_dir)
    
    os.makedirs(output_dir, exist_ok=True)
    assets_dir = os.path.join(output_dir, "assets")
    os.makedirs(assets_dir, exist_ok=True)
    
    # 写入 CSS
    with open(os.path.join(assets_dir, "style.css"), "w", encoding="utf-8") as f:
        f.write(generate_css())
    
    # 生成首页
    with open(os.path.join(output_dir, "index.html"), "w", encoding="utf-8") as f:
        f.write(generate_index(CHAPTERS))
    print("✅ index.html")
    
    # 生成每个章节页面
    for chapter in CHAPTERS:
        files = read_go_files(root_dir, chapter["dir"])
        if not files:
            print(f"⚠️  跳过 {chapter['dir']} (无 .go 文件)")
            continue
        
        page_html = generate_html_page(chapter, files, CHAPTERS, root_dir)
        out_path = os.path.join(output_dir, f"{chapter['dir']}.html")
        with open(out_path, "w", encoding="utf-8") as f:
            f.write(page_html)
        print(f"✅ {chapter['dir']}.html ({len(files)} files)")
    
    print(f"\n🎉 站点已生成到 {output_dir}/")


if __name__ == "__main__":
    root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    output = os.path.join(root, "site")
    build_site(root, output)
