import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import React, { useEffect, useState } from 'react';

const Features = [
  {
    title: 'Instant Management',
    description: 'Detect, monitor, and connect to Android devices instantly over USB or Wi-Fi.',
  },
  {
    title: 'Performance',
    description: 'Track CPU, Memory, and Network usage with real-time graphs right from your terminal.',
  },
  {
    title: 'App & Files',
    description: 'Browse files, install/uninstall apps, and pull/push data seamlessly with intuitive keybindings.',
  },
];

function Feature({title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="feature-card">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

function TerminalAnimation() {
  const [text, setText] = useState('');
  const fullText = '$ brew install SakshhamTheCoder/tap/adbt\n$ adbt';

  useEffect(() => {
    let currentText = '';
    let i = 0;
    const interval = setInterval(() => {
      currentText += fullText.charAt(i);
      setText(currentText);
      i++;
      if (i === fullText.length) clearInterval(interval);
    }, 40);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="terminal-mockup">
      <div className="terminal-header">
        <div className="terminal-dot"></div>
        <div className="terminal-dot"></div>
        <div className="terminal-dot"></div>
      </div>
      <div className="terminal-body">
        <pre style={{ margin: 0, background: 'transparent', padding: 0, whiteSpace: 'pre-wrap' }}>
          <code>{text}<span style={{ animation: 'blink 1s step-end infinite' }}>_</span></code>
        </pre>
      </div>
    </div>
  );
}

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero', 'heroBanner')}>
      <div className="container">
        <div className="row">
          <div className="col col--12" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
            <Heading as="h1" className="hero__title">
              {siteConfig.title}
            </Heading>
            <p className="hero__subtitle">{siteConfig.tagline}</p>
            
            <div className="hero-buttons-container">
              <Link
                className="adbt-button adbt-button-primary"
                to="/docs/installation">
                Get Started
              </Link>
              <Link
                className="adbt-button adbt-button-secondary"
                to="https://github.com/SakshhamTheCoder/adbt">
                View on GitHub
              </Link>
            </div>
          </div>
        </div>
        <TerminalAnimation />
      </div>
    </header>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} — Android Debug Bridge TUI`}
      description="Modern, keyboard-driven Terminal User Interface for ADB">
      <HomepageHeader />
      <main style={{ padding: '8rem 0' }}>
        <div className="container">
          {Features && Features.length > 0 && (
            <div className="row" style={{ rowGap: '5rem' }}>
              {Features.map((props, idx) => (
                <Feature key={idx} {...props} />
              ))}
            </div>
          )}
        </div>
      </main>
    </Layout>
  );
}
