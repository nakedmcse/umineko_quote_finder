(async () => {
    const API_BASE = '/api/v1';

    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    const audioIdInput = document.getElementById('audioIdInput');
    const audioIdBtn = document.getElementById('audioIdBtn');
    const randomBtn = document.getElementById('randomBtn');
    const clearBtn = document.getElementById('clearBtn');
    const characterSelect = document.getElementById('characterSelect');
    const episodeSelect = document.getElementById('episodeSelect');
    const browseBtn = document.getElementById('browseBtn');
    const resultsContainer = document.getElementById('resultsContainer');
    const langBtns = document.querySelectorAll('.lang-btn');

    let currentLang = 'en';
    let currentAudioId = null;
    let browseMode = false;
    let browseCharacter = '';
    let browseEpisode = 0;
    let browseOffset = 0;
    let browseTotal = 0;

    function createButterflies() {
        const container = document.getElementById('butterflies');
        const count = 8;

        for (let i = 0; i < count; i++) {
            const butterfly = document.createElement('div');
            butterfly.className = 'butterfly';
            butterfly.innerHTML = '🦋';
            butterfly.style.setProperty('--start-x', `${Math.random() * 100}vw`);
            butterfly.style.setProperty('--duration', `${20 + Math.random() * 15}s`);
            butterfly.style.setProperty('--delay', `${Math.random() * 20}s`);
            butterfly.style.left = `${Math.random() * 100}%`;
            container.appendChild(butterfly);
        }
    }

    async function loadCharacters() {
        try {
            const response = await fetch(`${API_BASE}/characters`);
            const characters = await response.json();

            const sorted = Object.entries(characters).sort((a, b) => a[1].localeCompare(b[1]));

            for (const [id, name] of sorted) {
                const option = document.createElement('option');
                option.value = id;
                option.textContent = name;
                characterSelect.appendChild(option);
            }
        } catch (error) {
            console.error('Failed to load characters:', error);
        }
    }

    function showLoading() {
        resultsContainer.innerHTML = `
            <div class="loading">
                <div class="loading-spinner"></div>
                <p class="loading-text">Searching through the fragments...</p>
            </div>
        `;
    }

    function showEmpty(message = 'No quotes found in this fragment.') {
        resultsContainer.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">🦋</div>
                <h3 class="empty-title">The Golden Land remains silent</h3>
                <p class="empty-subtitle">${message}</p>
            </div>
        `;
    }

    let currentQuery = '';
    let currentOffset = 0;
    let currentTotal = 0;
    const PAGE_SIZE = 30;

    function renderQuotes(results, query, total, offset) {
        stopAudio();
        currentAudioId = null;
        if (!results || results.length === 0) {
            showEmpty();
            return;
        }

        const start = offset + 1;
        const end = offset + results.length;

        const header = query ? `
            <div class="results-header">
                <span class="results-count">Showing <span>${start}-${end}</span> of <span>${total}</span> fragments for "${escapeHtml(query)}"</span>
            </div>
        ` : '';

        const quotes = results.map((item, index) => {
            const quote = item.quote || item;
            return `
                <article class="quote-card" style="--index: ${index}">
                    <span class="quote-mark">"</span>
                    <p class="quote-text">${quote.textHtml || escapeHtml(quote.text)}</p>
                    <div class="quote-meta">
                        <span class="quote-character">— ${escapeHtml(quote.character)}</span>
                        <div class="quote-details">
                            ${quote.episode ? `<span class="quote-episode">Episode ${quote.episode}</span>` : ''}

                            ${langToggleHTML(quote.audioId)}
                        </div>
                    </div>
                    ${audioPlayerHTML(quote.audioId, quote.characterId)}
                </article>
            `;
        }).join('');

        const totalPages = Math.ceil(total / PAGE_SIZE);
        const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

        const pagination = total > PAGE_SIZE ? `
            <div class="pagination">
                <button class="pagination-btn" id="prevPage" ${offset === 0 ? 'disabled' : ''}>◀ Previous</button>
                <span class="pagination-info">Page <span>${currentPage}</span> of <span>${totalPages}</span></span>
                <button class="pagination-btn" id="nextPage" ${end >= total ? 'disabled' : ''}>Next ▶</button>
            </div>
        ` : '';

        resultsContainer.innerHTML = header + `<div class="quotes-grid">${quotes}</div>` + pagination;

        if (total > PAGE_SIZE) {
            document.getElementById('prevPage')?.addEventListener('click', () => {
                if (currentOffset >= PAGE_SIZE) {
                    currentOffset -= PAGE_SIZE;
                    search(currentQuery, currentOffset);
                }
            });
            document.getElementById('nextPage')?.addEventListener('click', () => {
                if (currentOffset + PAGE_SIZE < currentTotal) {
                    currentOffset += PAGE_SIZE;
                    search(currentQuery, currentOffset);
                }
            });
        }
    }

    function renderFeaturedQuote(quote) {
        stopAudio();
        currentAudioId = quote.audioId ? quote.audioId.split(', ')[0] : null;
        resultsContainer.innerHTML = `
            <article class="featured-quote">
                <div class="featured-label">✦ A Fragment from the Sea ✦</div>
                <p class="featured-text">"${quote.textHtml || escapeHtml(quote.text)}"</p>
                <p class="featured-character">— ${escapeHtml(quote.character)}</p>
                ${quote.episode ? `<p class="featured-episode">Episode ${quote.episode}</p>` : ''}
                ${audioPlayerHTML(quote.audioId, quote.characterId)}
                ${langToggleHTML(quote.audioId)}
            </article>
        `;
    }

    function formatAudioIds(audioId) {
        return audioId.split(', ').map(id => id + '.ogg').join(', ');
    }

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    function langToggleHTML(audioId) {
        if (!audioId) return '';
        const firstId = escapeHtml(audioId.split(', ')[0]);
        const enActive = currentLang === 'en' ? ' active' : '';
        const jaActive = currentLang === 'ja' ? ' active' : '';
        return `<span class="lang-card-toggle" data-audio-id="${firstId}"><button class="lang-card-btn${enActive}" data-lang="en">EN</button><button class="lang-card-btn${jaActive}" data-lang="ja">JA</button></span>`;
    }

    function audioPlayerHTML(audioId, characterId) {
        if (!audioId) return '';
        const charId = escapeHtml(characterId || '');
        const ids = audioId.split(', ');
        let individualClips = '';
        for (let i = 0; i < ids.length; i++) {
            const eid = escapeHtml(ids[i]);
            individualClips += `<button class="audio-clip-btn" data-audio-id="${eid}" data-char-id="${charId}">&#9654; ${eid}.ogg</button>`;
        }
        let clipsHTML;
        if (ids.length > 1) {
            const allIds = ids.map(id => escapeHtml(id)).join(',');
            clipsHTML = `<div class="audio-clips"><button class="audio-clip-btn audio-combined-btn" data-char-id="${charId}" data-audio-ids="${allIds}">&#9654; Combined (${ids.length} clips)</button><button class="audio-expand-btn">&#9662; Individual</button></div><div class="audio-individual-clips">${individualClips}</div>`;
        } else {
            clipsHTML = `<div class="audio-clips">${individualClips}</div>`;
        }
        const savedVolume = localStorage.getItem('uminekoVolume') ?? "1.0";
        const volume = Math.floor(parseFloat(savedVolume) * 100);
        return `<div class="audio-player">
            ${clipsHTML}
            <div class="audio-controls">
                <div class="audio-track"><div class="audio-progress"></div></div>
                <div class="audio-volume"><span>VOL</span><input class="audio-volume-slider" type="range" min="0" max="100" step="1" value="${volume}"</div>
                <span class="audio-time">0:00 / 0:00</span>
            </div>
        </div>`;
    }

    let activeAudio = null;
    let activeBtn = null;
    let activePlayer = null;

    function formatTime(sec) {
        if (!sec || !isFinite(sec)) return '0:00';
        const m = Math.floor(sec / 60);
        const s = Math.floor(sec % 60);
        return m + ':' + (s < 10 ? '0' : '') + s;
    }

    function stopAudio() {
        if (activeAudio) {
            activeAudio.pause();
            activeAudio.removeAttribute('src');
            activeAudio.load();
        }
        if (activeBtn) {
            activeBtn.classList.remove('active');
        }
        if (activePlayer) {
            activePlayer.classList.remove('playing');
            const controls = activePlayer.querySelector('.audio-controls');
            if (controls) controls.classList.remove('visible');
            const progress = activePlayer.querySelector('.audio-progress');
            if (progress) progress.style.width = '0%';
            const timeEl = activePlayer.querySelector('.audio-time');
            if (timeEl) timeEl.textContent = '0:00 / 0:00';
        }
        activeBtn = null;
        activePlayer = null;
    }

    function setVolume(volume) {
        if (!activePlayer) return;
        activeAudio.volume = volume;
        localStorage.setItem('uminekoVolume', volume.toString());
    }

    function playAudioUrl(url, btn) {
        const player = btn.closest('.audio-player');
        if (activeAudio && activeBtn === btn) {
            if (activeAudio.paused) {
                activeAudio.play();
                player.classList.add('playing');
            } else {
                activeAudio.pause();
                player.classList.remove('playing');
            }
            return;
        }

        stopAudio();

        if (!activeAudio) {
            activeAudio = new Audio();
            const savedVolume = localStorage.getItem('uminekoVolume');
            if (savedVolume) {
                activeAudio.volume = parseFloat(savedVolume);
            }
            if (player) {
                const volumeSelector = player.querySelector('.audio-volume-slider');
                if (volumeSelector) {
                    volumeSelector.oninput = null;
                    volumeSelector.addEventListener('input', () => {
                        const v = parseFloat(volumeSelector.value) / 100;
                        setVolume(v);
                    })
                }
            }
            activeAudio.addEventListener('timeupdate', () => {
                if (!activePlayer) return;
                const progress = activePlayer.querySelector('.audio-progress');
                const timeEl = activePlayer.querySelector('.audio-time');
                if (progress && activeAudio.duration) {
                    progress.style.width = (activeAudio.currentTime / activeAudio.duration * 100) + '%';
                }
                if (timeEl) {
                    timeEl.textContent = formatTime(activeAudio.currentTime) + ' / ' + formatTime(activeAudio.duration);
                }
            });
            activeAudio.addEventListener('ended', () => {
                stopAudio();
            });
            activeAudio.addEventListener('loadedmetadata', () => {
                if (activePlayer) {
                    const timeEl = activePlayer.querySelector('.audio-time');
                    if (timeEl) {
                        timeEl.textContent = '0:00 / ' + formatTime(activeAudio.duration);
                    }
                }
            });
        }

        activeBtn = btn;
        activePlayer = player;
        btn.classList.add('active');
        const controls = player.querySelector('.audio-controls');
        if (controls) controls.classList.add('visible');
        activeAudio.src = url;
        activeAudio.play();
        player.classList.add('playing');
    }

    function playAudio(charId, audioId, btn) {
        playAudioUrl(`${API_BASE}/audio/${charId}/${audioId}`, btn);
    }

    resultsContainer.addEventListener('click', (e) => {
        const expandBtn = e.target.closest('.audio-expand-btn');
        if (expandBtn) {
            const player = expandBtn.closest('.audio-player');
            const individual = player.querySelector('.audio-individual-clips');
            if (individual) {
                const open = individual.classList.toggle('visible');
                expandBtn.innerHTML = open ? '&#9652; Individual' : '&#9662; Individual';
            }
            return;
        }
        const combinedBtn = e.target.closest('.audio-combined-btn');
        if (combinedBtn) {
            const charId = combinedBtn.dataset.charId;
            const audioIds = combinedBtn.dataset.audioIds;
            playAudioUrl(`${API_BASE}/audio/${charId}/combined?ids=${audioIds}`, combinedBtn);
            return;
        }
        const clipBtn = e.target.closest('.audio-clip-btn');
        if (clipBtn) {
            playAudio(clipBtn.dataset.charId, clipBtn.dataset.audioId, clipBtn);
            return;
        }
        const track = e.target.closest('.audio-track');
        if (track && activeAudio && activeAudio.duration) {
            const player = track.closest('.audio-player');
            if (player === activePlayer) {
                const rect = track.getBoundingClientRect();
                const ratio = (e.clientX - rect.left) / rect.width;
                activeAudio.currentTime = ratio * activeAudio.duration;
            }
        }
    });

    async function search(query, offset = 0) {
        if (!query.trim()) {
            showEmpty('Enter a search term to find quotes.');
            return;
        }

        browseMode = false;
        showLoading();

        try {
            const characterId = characterSelect.value;
            const episode = episodeSelect.value;
            let url = `${API_BASE}/search?q=${encodeURIComponent(query)}&limit=${PAGE_SIZE}&offset=${offset}&lang=${currentLang}`;
            if (characterId) {
                url += `&character=${characterId}`;
            }
            if (episode && episode !== '0') {
                url += `&episode=${episode}`;
            }

            const response = await fetch(url);
            const data = await response.json();

            currentQuery = query;
            currentOffset = data.offset;
            currentTotal = data.total;

            const results = data.results || [];

            renderQuotes(results, query, data.total, data.offset);
            updateURL();
        } catch (error) {
            console.error('Search failed:', error);
            showEmpty('Failed to search. Please try again.');
        }
    }

    async function getRandomQuote() {
        browseMode = false;
        showLoading();

        try {
            const characterId = characterSelect.value;
            const episode = episodeSelect.value;
            let url = `${API_BASE}/random?lang=${currentLang}`;
            if (characterId) {
                url += `&character=${characterId}`;
            }
            if (episode && episode !== '0') {
                url += `&episode=${episode}`;
            }
            const response = await fetch(url);
            const quote = await response.json();
            if (quote.error) {
                showEmpty('No quotes found for this character.');
                return;
            }
            renderFeaturedQuote(quote);
            updateURL();
        } catch (error) {
            console.error('Failed to get random quote:', error);
            showEmpty('Failed to retrieve a quote.');
        }
    }

    async function getQuoteByAudioId(audioId) {
        showLoading();

        try {
            const response = await fetch(`${API_BASE}/quote/${audioId}?lang=${currentLang}`);
            const quote = await response.json();
            if (quote.error) {
                showEmpty(`No quote found for audio ID "${escapeHtml(audioId)}".`);
                return;
            }
            renderFeaturedQuote(quote);
            updateURL();
        } catch (error) {
            console.error('Failed to get quote by audioId:', error);
            showEmpty('Failed to look up audio ID. Please try again.');
        }
    }

    searchBtn.addEventListener('click', () => {
        currentOffset = 0;
        search(searchInput.value, 0);
    });
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            currentOffset = 0;
            search(searchInput.value, 0);
        }
    });

    function lookupAudioId() {
        const id = audioIdInput.value.trim();
        if (!id) return;
        getQuoteByAudioId(id);
    }
    audioIdBtn.addEventListener('click', lookupAudioId);
    audioIdInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') lookupAudioId();
    });

    randomBtn.addEventListener('click', getRandomQuote);

    clearBtn.addEventListener('click', () => {
        stopAudio();
        searchInput.value = '';
        audioIdInput.value = '';
        characterSelect.value = '';
        episodeSelect.value = '0';
        resultsContainer.innerHTML = '';
        browseMode = false;
        browseBtn.disabled = true;
        updateURL();
    });

    characterSelect.addEventListener('change', () => {
        browseBtn.disabled = !characterSelect.value;
        if (!browseMode) {
            currentOffset = 0;
            if (searchInput.value.trim()) {
                search(searchInput.value, 0);
            }
        }
    });

    episodeSelect.addEventListener('change', () => {
        if (browseMode && browseCharacter) {
            browseEpisode = parseInt(episodeSelect.value) || 0;
            browseOffset = 0;
            browseCharacterDialogue(browseCharacter, browseOffset, browseEpisode);
        } else if (searchInput.value.trim()) {
            currentOffset = 0;
            search(searchInput.value, 0);
        }
    });

    browseBtn.addEventListener('click', () => {
        const characterId = characterSelect.value;
        if (!characterId) return;

        browseMode = true;
        browseCharacter = characterId;
        browseEpisode = parseInt(episodeSelect.value) || 0;
        browseOffset = 0;
        searchInput.value = '';
        browseCharacterDialogue(characterId, 0, browseEpisode);
    });

    resultsContainer.addEventListener('click', async (e) => {
        const btn = e.target.closest('.lang-card-btn');
        if (!btn || btn.classList.contains('active')) return;

        const toggle = btn.closest('.lang-card-toggle');
        const audioId = toggle.dataset.audioId;
        const newLang = btn.dataset.lang;

        for (const b of toggle.querySelectorAll('.lang-card-btn')) {
            b.disabled = true;
        }
        try {
            const response = await fetch(`${API_BASE}/quote/${audioId}?lang=${newLang}`);
            const quote = await response.json();
            if (quote.error) return;

            const card = btn.closest('.quote-card') || btn.closest('.featured-quote');
            const textEl = card.querySelector('.quote-text') || card.querySelector('.featured-text');
            if (textEl) {
                const isFeatured = textEl.classList.contains('featured-text');
                textEl.innerHTML = isFeatured
                    ? `"${quote.textHtml || escapeHtml(quote.text)}"`
                    : (quote.textHtml || escapeHtml(quote.text));
            }

            for (const b of toggle.querySelectorAll('.lang-card-btn')) {
                b.classList.toggle('active', b.dataset.lang === newLang);
            }
        } catch (error) {
            console.error('Failed to toggle language:', error);
        } finally {
            for (const b of toggle.querySelectorAll('.lang-card-btn')) {
                b.disabled = false;
            }
        }
    });

    async function browseCharacterDialogue(characterId, offset = 0, episode = 0) {
        showLoading();

        try {
            let url = `${API_BASE}/character/${characterId}?limit=${PAGE_SIZE}&offset=${offset}&lang=${currentLang}`;
            if (episode > 0) {
                url += `&episode=${episode}`;
            }

            const response = await fetch(url);
            const data = await response.json();

            browseOffset = data.offset;
            browseTotal = data.total;

            renderBrowseResults(data);
            updateURL();
        } catch (error) {
            console.error('Browse failed:', error);
            showEmpty('Failed to load character dialogue.');
        }
    }

    function renderBrowseResults(data) {
        stopAudio();
        currentAudioId = null;
        if (!data.quotes || data.quotes.length === 0) {
            showEmpty('No dialogue found for this character.');
            return;
        }

        const episodeLabel = browseEpisode > 0 ? ` — Episode ${browseEpisode}` : '';
        const header = `
            <div class="browse-header">
                <h2 class="browse-title">${escapeHtml(data.character)}${episodeLabel}</h2>
                <p class="browse-subtitle">Showing lines ${data.offset + 1}-${data.offset + data.quotes.length} of ${data.total} in story order</p>
            </div>
        `;

        const quotes = data.quotes.map((quote, index) => {
            const lineNum = data.offset + index + 1;
            return `
                <article class="quote-card" style="--index: ${index}">
                    <span class="quote-number">#${lineNum}</span>
                    <span class="quote-mark">"</span>
                    <p class="quote-text">${quote.textHtml || escapeHtml(quote.text)}</p>
                    <div class="quote-meta">
                        <span class="quote-character">— ${escapeHtml(quote.character)}</span>
                        <div class="quote-details">
                            ${quote.episode ? `<span class="quote-episode">Episode ${quote.episode}</span>` : ''}

                            ${langToggleHTML(quote.audioId)}
                        </div>
                    </div>
                    ${audioPlayerHTML(quote.audioId, quote.characterId)}
                </article>
            `;
        }).join('');

        const totalPages = Math.ceil(data.total / PAGE_SIZE);
        const currentPage = Math.floor(data.offset / PAGE_SIZE) + 1;

        const pagination = data.total > PAGE_SIZE ? `
            <div class="pagination">
                <button class="pagination-btn" id="browsePrev" ${data.offset === 0 ? 'disabled' : ''}>◀ Previous</button>
                <span class="pagination-info">Page <span>${currentPage}</span> of <span>${totalPages}</span></span>
                <button class="pagination-btn" id="browseNext" ${data.offset + data.quotes.length >= data.total ? 'disabled' : ''}>Next ▶</button>
            </div>
        ` : '';

        resultsContainer.innerHTML = header + `<div class="quotes-grid">${quotes}</div>` + pagination;

        if (data.total > PAGE_SIZE) {
            document.getElementById('browsePrev')?.addEventListener('click', () => {
                if (browseOffset >= PAGE_SIZE) {
                    browseOffset -= PAGE_SIZE;
                    browseCharacterDialogue(browseCharacter, browseOffset, browseEpisode);
                }
            });
            document.getElementById('browseNext')?.addEventListener('click', () => {
                if (browseOffset + PAGE_SIZE < browseTotal) {
                    browseOffset += PAGE_SIZE;
                    browseCharacterDialogue(browseCharacter, browseOffset, browseEpisode);
                }
            });
        }
    }

    function updateURL() {
        const params = new URLSearchParams();

        if (browseMode && browseCharacter) {
            params.set('browse', browseCharacter);
        } else if (searchInput.value.trim()) {
            params.set('q', searchInput.value.trim());
            if (characterSelect.value) {
                params.set('character', characterSelect.value);
            }
        }

        const episode = episodeSelect.value;
        if (episode && episode !== '0') {
            params.set('episode', episode);
        }

        const offset = browseMode ? browseOffset : currentOffset;
        if (offset > 0) {
            params.set('offset', offset);
        }

        if (currentLang !== 'en') {
            params.set('lang', currentLang);
        }

        const qs = params.toString();
        const newURL = qs ? `?${qs}` : window.location.pathname;
        history.pushState(null, '', newURL);
    }

    function loadFromURL() {
        const params = new URLSearchParams(window.location.search);

        const lang = params.get('lang') || 'en';
        currentLang = lang;
        for (const b of langBtns) {
            b.classList.toggle('active', b.dataset.lang === lang);
        }

        const episode = params.get('episode') || '0';
        episodeSelect.value = episode;

        const browse = params.get('browse');
        const q = params.get('q');
        const offset = parseInt(params.get('offset')) || 0;

        if (browse) {
            characterSelect.value = browse;
            browseBtn.disabled = false;
            browseMode = true;
            browseCharacter = browse;
            browseEpisode = parseInt(episode) || 0;
            browseOffset = offset;
            searchInput.value = '';
            browseCharacterDialogue(browse, offset, browseEpisode);
        } else if (q) {
            searchInput.value = q;
            const character = params.get('character') || '';
            characterSelect.value = character;
            browseBtn.disabled = !character;
            currentOffset = offset;
            search(q, offset);
        } else {
            characterSelect.value = '';
            browseBtn.disabled = true;
            searchInput.value = '';
            getRandomQuote();
        }
    }

    window.addEventListener('popstate', loadFromURL);

    for (const btn of langBtns) {
        btn.addEventListener('click', () => {
            const newLang = btn.dataset.lang;
            if (newLang === currentLang) return;

            currentLang = newLang;
            for (const b of langBtns) {
                b.classList.toggle('active', b.dataset.lang === newLang);
            }

            if (browseMode && browseCharacter) {
                browseCharacterDialogue(browseCharacter, browseOffset, browseEpisode);
            } else if (searchInput.value.trim()) {
                search(searchInput.value, currentOffset);
            } else if (currentAudioId) {
                getQuoteByAudioId(currentAudioId);
            } else {
                getRandomQuote();
            }
        });
    }

    createButterflies();
    await loadCharacters();
    loadFromURL();
})();
