import React, { useState, useEffect } from 'react';
import { Capacitor } from '@capacitor/core';
import { App as CapacitorApp } from '@capacitor/app';
import { Haptics, ImpactStyle } from '@capacitor/haptics';
import { StatusBar, Style } from '@capacitor/status-bar';
import './App.css';

interface AppInfo {
  name: string;
  id: string;
  build: string;
  version: string;
}

function App() {
  const [appInfo, setAppInfo] = useState<AppInfo | null>(null);
  const [platform, setPlatform] = useState<string>('');
  const [isNative, setIsNative] = useState<boolean>(false);

  useEffect(() => {
    initializeApp();
  }, []);

  const initializeApp = async () => {
    // Get platform information
    setPlatform(Capacitor.getPlatform());
    setIsNative(Capacitor.isNativePlatform());

    // Get app information
    if (Capacitor.isNativePlatform()) {
      try {
        const info = await CapacitorApp.getInfo();
        setAppInfo(info);
      } catch (error) {
        console.error('Error getting app info:', error);
      }
    }
  };

  const triggerHaptic = async (style: ImpactStyle) => {
    if (Capacitor.isNativePlatform()) {
      try {
        await Haptics.impact({ style });
      } catch (error) {
        console.error('Error triggering haptic:', error);
      }
    } else {
      // Fallback for web
      console.log(`Haptic feedback: ${style}`);
    }
  };

  const toggleStatusBar = async () => {
    if (Capacitor.isNativePlatform()) {
      try {
        await StatusBar.setStyle({ style: Style.Dark });
        console.log('Status bar style changed to dark');
      } catch (error) {
        console.error('Error changing status bar:', error);
      }
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>🚀 React + Capacitor</h1>
        
        <div className="platform-info">
          <h2>Platform Information</h2>
          <p><strong>Platform:</strong> {platform}</p>
          <p><strong>Is Native:</strong> {isNative ? 'Yes' : 'No'}</p>
          
          {appInfo && (
            <div>
              <p><strong>App Name:</strong> {appInfo.name}</p>
              <p><strong>App ID:</strong> {appInfo.id}</p>
              <p><strong>Version:</strong> {appInfo.version}</p>
              <p><strong>Build:</strong> {appInfo.build}</p>
            </div>
          )}
        </div>

        <div className="controls">
          <h2>Native Features</h2>
          
          <div className="button-group">
            <button 
              onClick={() => triggerHaptic(ImpactStyle.Light)}
              className="haptic-button light"
            >
              Light Haptic
            </button>
            <button 
              onClick={() => triggerHaptic(ImpactStyle.Medium)}
              className="haptic-button medium"
            >
              Medium Haptic
            </button>
            <button 
              onClick={() => triggerHaptic(ImpactStyle.Heavy)}
              className="haptic-button heavy"
            >
              Heavy Haptic
            </button>
          </div>

          <button 
            onClick={toggleStatusBar}
            className="status-bar-button"
            disabled={!isNative}
          >
            Toggle Status Bar (Dark)
          </button>
        </div>

        <div className="instructions">
          <h3>Next Steps:</h3>
          <ul>
            <li>Run <code>npm start</code> for web development</li>
            <li>Run <code>npx cap run ios</code> to test on iOS</li>
            <li>Run <code>npx cap run android</code> to test on Android</li>
            <li>Use <code>npx cap sync</code> after code changes</li>
          </ul>
        </div>
      </header>
    </div>
  );
}

export default App;
