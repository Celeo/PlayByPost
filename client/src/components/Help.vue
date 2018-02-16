<template lang="pug">
  div
    v-container(grid-list-md)
      v-layout(row)
        v-flex(xs8 offset-xs2)
          h2.mb-4 How to use this site
          p.
            Welcome to this <strong>Play by Post</strong> web app. It's used to play RPGs in a <i>post-by-post</i> format: each
            player posts what their character is doing as they have time during their day, and their story unfolds asynchronously,
            without all players having to meet physically or virtually at the same time to play.
          p.mt-4.
            This app is designed to facilitate that effort: user's posts are layed out in a simple unnested old-school forum layout,
            one above the other, with the most recent posts at the bottom.
          p.mt-4.
            As this board is intended to be viewed by anyone but only allow participating users to post, anyone can view the existing
            posts without having to have an account or log in, but only users with accounts can make additional posts to the board.
          h2 Layout
          p.mt-4 The layout of the site is as follows:
          ul
            li.
              At the top of the page is the <strong>toolbar</strong>. It contains, left to right, the sidenav button, the web app title,
              a blank area, and on the far right, an indicator telling you which character you are logged in as (if you're logged in).
            li.
              On the left side of the page is nothing special, until you click the <i>sidenav button</i> (also known as a <i>hamburger button</i>)
              on the toolbar. Once you do, the side navigation panel will slide out from the left of the page and show additional links:
            ul.ml-5
              li <strong>Posts</strong> - view the posts on the board (basically, go to the home page)
              li <strong>Help</strong> - view this page
              li <strong>Log in</strong> - log into your character/user account
              li <strong>Register</strong> - create a new account
              li <strong>Change password</strong> - change your password
              li <strong>Log out</strong> - log out of your account
              li <strong>Close this</strong> - close the side navigation panel
            span You can also close the panel by clicking or tapping anywhere on the shaded area that appears
            li.
              In the center of the screen lies that page's content, whether that's posts, the login form, this text, etc.
          h2.mt-4 A note on passwords
          p.mt-4.
            With all of the security concerns and leaks and such going on these days, I've made sure that I'm using industry best-practices
            for hashing and salting your password. I never see it. The database content for your password will be something unrecognizable
            and irreversible, like <strong>$2a$10$rBV2JDeWW3.vKyeQcM8fFO4777l4bVeQgDL6VIkxqlzQ7TCalQvla</strong>. If you have a software
            background or are just curious, this app is <i>open source</i> and can be viewed on
            <a href="https://github.com/Celeo/PlayByPost">my GitHub account</a>.
          h2.mt-4.mb-4 Posting
          p.
            New posts can be made by logged-in users by going to the <strong>Posts</strong> page and scrolling to the bottom. There, you will
            an area titled <strong>New post</strong>. The editor below that is a <i>WYSIWYG</i> editor: "What you see if what you get." It's
            more like a word processing program's editor than what you can see on some other play by post forums, like Paizo's.
          p.
            Type your post into the text area. You can use the buttons on the top of the editor to bold your text, add italics, add lists, and
            change the color to blue to denote "out of character" text. The rightmost button will remove all text formatting from whatever you
            have highlighted when you click it. Additionally, some common <i>keyboard shortcuts</i> work, like <strong>Ctrl+B</strong> to bold.
          p.
            For now, rolling dice works as it does on Paizo's play by post forums - manually typing in text in the form of
            <strong>[dice=<code>action you're doing</code>] <code>dice rolls and modifiers</code> [/dice]</strong>.
          h4.mb-1 For example:
          p.
            If you were rolling for <strong>Perception</strong> and have a modifier of <strong>+3</strong>, then you would type in
            <code>[dice=Perception]1d20+3[/dice]</code>. If you're rolling damage for a weapon attack, you might type in
            <code>[dice=Damage]2d4-1[/dice]</code>.
          p.
            For now, you'll know what your dice values are <strong>after you submit your post</strong>. Because this is a little different
            than Paizo's forum, you may need to roll multiple dice in the same post, or make multiple posts as you learn what you got.
          h3.mt-2 Submitting
          p.
            Once you have written what your character does, click the blue <strong>Save</strong> button under the editor (it'll only be clickable
            if you've written something). You posting will be saved and shown on the board at the bottom.
          h2.mt-4.mb-4 Mobile
          p.
            This site should work on mobile. The site should recognize the viewing size of your device and shift the content and spacing accordingly.
          h2.mt-4.mb-4 Bugs / issues / requests
          p.
            You can submit issues and requests to <a href="mailto:celeodor@gmail.com">celeodor@gmail.com</a> or, if you use GitHub, the
            <a href="https://github.com/Celeo/PlayByPost">repo</a>.
</template>

<script>
export default {}
</script>

<style lang="stylus" scoped>
code
  background-color rgba(0, 0, 0, 0.05) !important
</style>
